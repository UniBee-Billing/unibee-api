package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	v1 "unibee-api/api/open/payment"
	redismqcmd "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/api"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/invoice/handler"
	"unibee-api/internal/logic/payment/callback"
	"unibee-api/internal/logic/payment/event"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
)

func GatewayPaymentCreate(ctx context.Context, createPayContext *ro.CreatePayContext) (gatewayInternalPayResult *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.Pay.BizType > 0, "pay bizType is nil")
	utility.Assert(createPayContext.Gateway != nil, "pay gateway is nil")
	utility.Assert(createPayContext.Pay != nil, "pay is nil")
	utility.Assert(len(createPayContext.Pay.BizId) > 0, "BizId Invalid")
	utility.Assert(createPayContext.Pay.GatewayId > 0, "pay gatewayId is nil")
	utility.Assert(createPayContext.Pay.TotalAmount > 0, "TotalAmount Invalid")
	//utility.Assert(len(createPayContext.Pay.CountryCode) > 0, "countryCode is nil")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(createPayContext.Pay.MerchantId > 0, "merchantId Invalid")
	utility.Assert(createPayContext.Pay.CompanyId > 0, "companyId Invalid")
	// 查询并处理所有待支付订单 todo mark

	createPayContext.Pay.Status = consts.TO_BE_PAID
	//createPayContext.Pay.AdditionalData = todo mark
	createPayContext.Pay.PaymentId = utility.CreatePaymentId()
	createPayContext.Pay.OpenApiId = createPayContext.OpenApiId
	if len(createPayContext.Invoice.InvoiceId) > 0 {
		createPayContext.Pay.InvoiceId = createPayContext.Invoice.InvoiceId
	}
	createPayContext.Pay.InvoiceData = utility.MarshalToJsonString(createPayContext.Invoice)
	if createPayContext.MediaData == nil {
		createPayContext.MediaData = make(map[string]string)
	}
	createPayContext.MediaData["BizType"] = strconv.Itoa(createPayContext.Pay.BizType)
	createPayContext.MediaData["PaymentId"] = createPayContext.Pay.PaymentId
	//toSave.setServiceRate(iMerchantInfoService.getServiceDeductPoint(toSave.getMerchantId(),toSave.getChannelId()));//记录当下的服务费率
	redisKey := fmt.Sprintf("createPay-merchantId:%d-bizId:%s", createPayContext.Pay.MerchantId, createPayContext.Pay.BizId)
	isDuplicatedInvoke := false

	//- 并发调用严重，加上Redis排它锁, todo mark database lock
	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		return nil, gerror.Newf(`too fast duplicate call %s`, createPayContext.Pay.BizId)
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCreated, createPayContext.Pay.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//transaction gateway refund
			createPayContext.Pay.UniqueId = createPayContext.Pay.PaymentId
			createPayContext.Pay.CreateTime = gtime.Now().Timestamp()
			insert, err := dao.Payment.Ctx(ctx).Data(createPayContext.Pay).OmitNil().Insert(createPayContext.Pay)
			if err != nil {
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				return err
			}
			createPayContext.Pay.Id = id

			//调用远端接口，这里的正向有坑，如果远端执行成功，事务却提交失败是无法回滚的 todo mark

			gatewayInternalPayResult, err = api.GetGatewayServiceProvider(ctx, createPayContext.Pay.GatewayId).GatewayPayment(ctx, createPayContext)
			if err != nil {
				return err
			}
			jsonData, err := gjson.Marshal(gatewayInternalPayResult)
			if err != nil {
				return err
			}
			var automatic = 0
			if gatewayInternalPayResult.Status == consts.PAY_SUCCESS && createPayContext.PayImmediate {
				automatic = 1
			}
			createPayContext.Pay.PaymentData = string(jsonData)
			createPayContext.Pay.Status = int(gatewayInternalPayResult.Status)
			createPayContext.Pay.GatewayPaymentId = gatewayInternalPayResult.GatewayPaymentId
			createPayContext.Pay.GatewayPaymentIntentId = gatewayInternalPayResult.GatewayPaymentIntentId
			gatewayInternalPayResult.PaymentId = createPayContext.Pay.PaymentId
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().PaymentData:            string(jsonData),
				dao.Payment.Columns().Automatic:              automatic,
				dao.Payment.Columns().Status:                 createPayContext.Pay.Status,
				dao.Payment.Columns().Link:                   gatewayInternalPayResult.Link,
				dao.Payment.Columns().GatewayPaymentId:       gatewayInternalPayResult.GatewayPaymentId,
				dao.Payment.Columns().GatewayPaymentIntentId: gatewayInternalPayResult.GatewayPaymentIntentId},
				g.Map{dao.Payment.Columns().Id: id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				return err
			}

			return nil
		})
		if err == nil {
			return redismq.CommitTransaction, nil
		} else {
			return redismq.RollbackTransaction, err
		}
	})

	// create new invoice, ignore errors
	invoice, _ := handler.CreateOrUpdateInvoiceFromPayment(ctx, createPayContext.Invoice, createPayContext.Pay)
	callback.GetPaymentCallbackServiceProvider(ctx, createPayContext.Pay.BizType).PaymentCreateCallback(ctx, createPayContext.Pay, invoice)
	if createPayContext.Pay.Status == consts.PAY_SUCCESS {
		callback.GetPaymentCallbackServiceProvider(ctx, createPayContext.Pay.BizType).PaymentSuccessCallback(ctx, createPayContext.Pay, invoice)
	}

	if err != nil {
		return nil, err
	} else {
		//交易事件记录
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     createPayContext.Pay.PaymentId,
			Fee:       createPayContext.Pay.TotalAmount,
			EventType: event.SentForSettle.Type,
			Event:     event.SentForSettle.Desc,
			OpenApiId: createPayContext.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", createPayContext.Pay.PaymentId, "SentForSettle"),
		})
	}
	return gatewayInternalPayResult, nil
}

func CreateSubInvoicePayment(ctx context.Context, sub *entity.Subscription, invoice *ro.InvoiceDetailSimplify, billingReason string) (gatewayInternalPayResult *ro.CreatePayInternalResp, err error) {
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	var mobile = ""
	var firstName = ""
	var lastName = ""
	var gender = ""
	var email = ""
	if user != nil {
		mobile = user.Mobile
		firstName = user.FirstName
		lastName = user.LastName
		gender = user.Gender
		email = user.Email
	}
	gateway := query.GetSubscriptionTypeGatewayById(ctx, sub.GatewayId)
	if gateway == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice gateway not found")
	}
	merchantInfo := query.GetMerchantInfoById(ctx, sub.MerchantId)
	if merchantInfo == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice merchantInfo not found")
	}
	return GatewayPaymentCreate(ctx, &ro.CreatePayContext{
		Gateway: gateway,
		Pay: &entity.Payment{
			SubscriptionId:  sub.SubscriptionId,
			BizId:           sub.SubscriptionId,
			BizType:         consts.BIZ_TYPE_SUBSCRIPTION,
			AuthorizeStatus: consts.AUTHORIZED,
			UserId:          sub.UserId,
			GatewayId:       int64(gateway.Id),
			TotalAmount:     invoice.TotalAmount,
			Currency:        invoice.Currency,
			CountryCode:     sub.CountryCode,
			MerchantId:      sub.MerchantId,
			CompanyId:       merchantInfo.CompanyId,
			BillingReason:   billingReason,
		},
		Platform:      "WEB",
		DeviceType:    "Web",
		ShopperUserId: strconv.FormatInt(sub.UserId, 10),
		ShopperEmail:  email,
		ShopperLocale: "en",
		Mobile:        mobile,
		Invoice:       invoice,
		ShopperName: &v1.OutShopperName{
			FirstName: firstName,
			LastName:  lastName,
			Gender:    gender,
		},
		MediaData:              map[string]string{"BillingReason": billingReason},
		MerchantOrderReference: sub.SubscriptionId,
		PayMethod:              1, //automatic
		DaysUtilDue:            5, // todo mark
		GatewayPaymentMethod:   sub.GatewayDefaultPaymentMethod,
	})
}
