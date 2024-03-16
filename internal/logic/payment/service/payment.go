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
	"strings"
	"unibee/api/bean"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/currency"
	email2 "unibee/internal/logic/email"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/logic/payment/event"
	handler2 "unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

func GatewayPaymentCreate(ctx context.Context, createPayContext *gateway_bean.GatewayNewPaymentReq) (gatewayInternalPayResult *gateway_bean.GatewayNewPaymentResp, err error) {
	utility.Assert(createPayContext.Pay.BizType > 0, "pay bizType is nil")
	utility.Assert(createPayContext.Gateway != nil, "pay gateway is nil")
	utility.Assert(createPayContext.Pay != nil, "pay is nil")
	utility.Assert(len(createPayContext.Pay.ExternalPaymentId) > 0, "BizId Invalid")
	utility.Assert(createPayContext.Pay.GatewayId > 0, "pay gatewayId is nil")
	utility.Assert(createPayContext.Pay.TotalAmount > 0, "TotalAmount Invalid")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(createPayContext.Pay.MerchantId > 0, "merchantId Invalid")
	utility.Assert(currency.IsFiatCurrencySupport(createPayContext.Pay.Currency), "currency not support")

	createPayContext.Pay.Status = consts.PaymentCreated
	createPayContext.Pay.PaymentId = utility.CreatePaymentId()
	createPayContext.Pay.InvoiceData = utility.MarshalToJsonString(createPayContext.Invoice)
	if createPayContext.Metadata == nil {
		createPayContext.Metadata = make(map[string]string)
	}
	createPayContext.Metadata["PaymentId"] = createPayContext.Pay.PaymentId
	createPayContext.Metadata["MerchantId"] = strconv.FormatUint(createPayContext.Pay.MerchantId, 10)
	createPayContext.Pay.MetaData = utility.MarshalToJsonString(createPayContext.Metadata)
	redisKey := fmt.Sprintf("createPay-merchantId:%d-externalPaymentId:%s", createPayContext.Pay.MerchantId, createPayContext.Pay.ExternalPaymentId)
	isDuplicatedInvoke := false

	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		return nil, gerror.Newf(`too fast duplicate call %s`, createPayContext.Pay.ExternalPaymentId)
	}

	if createPayContext.Gateway.GatewayType == consts.GatewayTypeCrypto {
		//crypto payment
		if len(createPayContext.Pay.GasPayer) > 0 {
			utility.Assert(strings.Contains("merchant|user", createPayContext.Pay.GasPayer), "crypto payment gasPayer should one of merchant|user")
		} else {
			createPayContext.Pay.GasPayer = "user" // default user pay the gas
		}
		trans, err := api.GetGatewayServiceProvider(ctx, createPayContext.Pay.GatewayId).GatewayCryptoFiatTrans(ctx, &gateway_bean.GatewayCryptoFromCurrencyAmountDetailReq{
			Amount:      createPayContext.Pay.TotalAmount,
			Currency:    createPayContext.Pay.Currency,
			CountryCode: createPayContext.Pay.CountryCode,
			Gateway:     createPayContext.Gateway,
		})
		if err != nil {
			return nil, err
		}
		createPayContext.Pay.CryptoAmount = trans.CryptoAmount
		createPayContext.Pay.CryptoCurrency = trans.CryptoCurrency
	}
	var invoice *entity.Invoice
	if createPayContext.Invoice != nil {
		invoice, err = handler.CreateOrUpdateInvoiceForNewPayment(ctx, createPayContext.Invoice, createPayContext.Pay)
		if err != nil {
			return nil, err
		}
		createPayContext.Pay.InvoiceId = invoice.InvoiceId
	}
	if createPayContext.DaysUtilDue == 0 {
		createPayContext.DaysUtilDue = 3 //default 3 days expire
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCreated, createPayContext.Pay.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//transaction gateway refund
			createPayContext.Pay.UniqueId = createPayContext.Pay.PaymentId
			createPayContext.Pay.CreateTime = gtime.Now().Timestamp()
			createPayContext.Pay.ExpireTime = createPayContext.Pay.CreateTime + int64(createPayContext.DaysUtilDue*86400)
			insert, err := dao.Payment.Ctx(ctx).Data(createPayContext.Pay).OmitNil().Insert(createPayContext.Pay)
			if err != nil {
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				return err
			}
			createPayContext.Pay.Id = id

			gatewayInternalPayResult, err = api.GetGatewayServiceProvider(ctx, createPayContext.Pay.GatewayId).GatewayNewPayment(ctx, createPayContext)
			if err != nil {
				return err
			}
			jsonData, err := gjson.Marshal(gatewayInternalPayResult)
			if err != nil {
				return err
			}
			var automatic = 0
			if gatewayInternalPayResult.Status == consts.PaymentSuccess && createPayContext.PayImmediate {
				automatic = 1
			}
			createPayContext.Pay.PaymentData = string(jsonData)
			createPayContext.Pay.Status = int(gatewayInternalPayResult.Status)
			createPayContext.Pay.GatewayPaymentId = gatewayInternalPayResult.GatewayPaymentId
			createPayContext.Pay.GatewayPaymentIntentId = gatewayInternalPayResult.GatewayPaymentIntentId
			gatewayInternalPayResult.PaymentId = createPayContext.Pay.PaymentId
			// unibee payment link
			paymentLink := link.GetPaymentLink(createPayContext.Pay.PaymentId)
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().PaymentData:            string(jsonData),
				dao.Payment.Columns().Automatic:              automatic,
				dao.Payment.Columns().Link:                   paymentLink,
				dao.Payment.Columns().GatewayLink:            gatewayInternalPayResult.Link,
				dao.Payment.Columns().GatewayPaymentId:       gatewayInternalPayResult.GatewayPaymentId,
				dao.Payment.Columns().GatewayPaymentIntentId: gatewayInternalPayResult.GatewayPaymentIntentId},
				g.Map{dao.Payment.Columns().Id: id, dao.Payment.Columns().Status: consts.PaymentCreated})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				return err
			}
			gatewayInternalPayResult.Link = paymentLink
			createPayContext.Pay.Link = paymentLink
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction, nil
		} else {
			return redismq.RollbackTransaction, err
		}
	})
	if err != nil {
		return nil, err
	}

	gatewayInternalPayResult.Invoice = invoice
	callback.GetPaymentCallbackServiceProvider(ctx, createPayContext.Pay.BizType).PaymentCreateCallback(ctx, createPayContext.Pay, gatewayInternalPayResult.Invoice)
	if createPayContext.Pay.Status == consts.PaymentSuccess {
		req := &handler2.HandlePayReq{
			PaymentId:              createPayContext.Pay.PaymentId,
			GatewayPaymentIntentId: gatewayInternalPayResult.GatewayPaymentIntentId,
			GatewayPaymentId:       gatewayInternalPayResult.GatewayPaymentId,
			GatewayPaymentMethod:   gatewayInternalPayResult.GatewayPaymentMethod,
			PayStatusEnum:          consts.PaymentSuccess,
			TotalAmount:            createPayContext.Pay.TotalAmount,
			PaymentAmount:          createPayContext.Pay.TotalAmount,
			PaidTime:               gtime.Now(),
		}
		err = handler2.HandlePaySuccess(ctx, req)
	}
	invoice, err = handler.CreateOrUpdateInvoiceForNewPayment(ctx, createPayContext.Invoice, createPayContext.Pay)
	gatewayInternalPayResult.Invoice = invoice
	if err != nil {
		return nil, err
	}

	event.SaveEvent(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     createPayContext.Pay.PaymentId,
		Fee:       createPayContext.Pay.TotalAmount,
		EventType: event.SentForSettle.Type,
		Event:     event.SentForSettle.Desc,
		UniqueNo:  fmt.Sprintf("%s_%s", createPayContext.Pay.PaymentId, "SentForSettle"),
	})
	return gatewayInternalPayResult, nil
}

func CreateSubInvoiceAutomaticPayment(ctx context.Context, sub *entity.Subscription, invoice *entity.Invoice) (gatewayInternalPayResult *gateway_bean.GatewayNewPaymentResp, err error) {
	user := query.GetUserAccountById(ctx, uint64(sub.UserId))
	var email = ""
	if user != nil {
		email = user.Email
	}
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	if gateway == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice gateway not found")
	}

	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	if merchant == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice merchantInfo not found")
	}
	res, err := GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		PayImmediate: true,
		Gateway:      gateway,
		Pay: &entity.Payment{
			SubscriptionId:    sub.SubscriptionId,
			ExternalPaymentId: sub.SubscriptionId,
			BizType:           consts.BizTypeSubscription,
			AuthorizeStatus:   consts.Authorized,
			UserId:            sub.UserId,
			GatewayId:         gateway.Id,
			TotalAmount:       invoice.TotalAmount,
			Currency:          invoice.Currency,
			CryptoAmount:      invoice.CryptoAmount,
			CryptoCurrency:    invoice.CryptoCurrency,
			CountryCode:       sub.CountryCode,
			MerchantId:        sub.MerchantId,
			CompanyId:         merchant.CompanyId,
			Automatic:         1,
			BillingReason:     invoice.InvoiceName,
			GasPayer:          sub.GasPayer,
		},
		ExternalUserId:       strconv.FormatInt(sub.UserId, 10),
		Email:                email,
		Invoice:              bean.SimplifyInvoice(invoice),
		Metadata:             map[string]string{"BillingReason": invoice.InvoiceName},
		GatewayPaymentMethod: sub.GatewayDefaultPaymentMethod,
	})
	if err == nil && res.Status != consts.PaymentSuccess {
		//need send invoice for authorised
		payment := query.GetPaymentByPaymentId(ctx, res.PaymentId)
		if payment != nil {
			oneUser := query.GetUserAccountById(ctx, uint64(sub.UserId))
			plan := query.GetPlanById(ctx, sub.PlanId)
			if plan != nil && oneUser != nil {
				err = email2.SendTemplateEmail(ctx, merchant.Id, oneUser.Email, oneUser.TimeZone, email2.TemplateSubscriptionNeedAuthorized, "", &email2.TemplateVariable{
					UserName:            oneUser.FirstName + " " + oneUser.LastName,
					MerchantProductName: plan.PlanName,
					MerchantCustomEmail: merchant.Email,
					MerchantName:        merchant.Name,
					PaymentAmount:       utility.ConvertCentToDollarStr(invoice.TotalAmount, invoice.Currency),
					Currency:            strings.ToUpper(invoice.Currency),
					PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
				})
				if err != nil {
					g.Log().Errorf(ctx, "CreateSubInvoiceAutomaticPayment SendTemplateEmail err:%s", err.Error())
				}
			}
		}
	}
	return res, err
}
