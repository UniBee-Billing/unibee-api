package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"unibee/api/bean"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/controller/link"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/discount"
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
	utility.Assert(len(createPayContext.Pay.ExternalPaymentId) > 0, "ExternalPaymentId Invalid")
	utility.Assert(createPayContext.Pay.GatewayId > 0, "pay gatewayId is nil")
	utility.Assert(createPayContext.Pay.TotalAmount > 0, "TotalAmount Invalid")
	utility.Assert(len(createPayContext.Pay.Currency) > 0, "currency is nil")
	utility.Assert(createPayContext.Pay.MerchantId > 0, "merchantId Invalid")
	utility.Assert(createPayContext.Invoice != nil, "invoice is nil")
	createPayContext.Pay.Currency = strings.ToUpper(createPayContext.Pay.Currency)
	createPayContext.Invoice.Currency = strings.ToUpper(createPayContext.Invoice.Currency)
	utility.Assert(currency.IsFiatCurrencySupport(createPayContext.Pay.Currency), "currency not support")

	createPayContext.Pay.Status = consts.PaymentCreated
	createPayContext.Pay.PaymentId = utility.CreatePaymentId()
	createPayContext.Pay.InvoiceData = utility.MarshalToJsonString(createPayContext.Invoice)
	if createPayContext.Metadata == nil {
		createPayContext.Metadata = make(map[string]interface{})
	}
	createPayContext.Metadata["PaymentId"] = createPayContext.Pay.PaymentId
	createPayContext.Metadata["MerchantId"] = strconv.FormatUint(createPayContext.Pay.MerchantId, 10)
	createPayContext.Pay.MetaData = utility.MarshalToJsonString(createPayContext.Metadata)
	redisKey := fmt.Sprintf("createPay-merchantId:%d-externalPaymentId:%s", createPayContext.Pay.MerchantId, createPayContext.Pay.ExternalPaymentId)
	isDuplicatedInvoke := false

	if createPayContext.Gateway.GatewayType == consts.GatewayTypeWireTransfer {
		utility.Assert(createPayContext.Pay.TotalAmount >= createPayContext.Gateway.MinimumAmount, "Total Amount not reach the gateway's minimum amount")
	}

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
	if createPayContext.DaysUtilDue == 0 {
		createPayContext.DaysUtilDue = 3 //default 3 days expire
	}

	var invoice *entity.Invoice
	if createPayContext.Invoice.Id > 0 {
		createPayContext.Pay.InvoiceId = createPayContext.Invoice.InvoiceId
	} else {
		createPayContext.Pay.InvoiceId = utility.CreateInvoiceId()
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCreated, createPayContext.Pay.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			//transaction gateway payment
			createPayContext.Pay.UniqueId = createPayContext.Pay.PaymentId
			createPayContext.Pay.CreateTime = gtime.Now().Timestamp()
			createPayContext.Pay.ExpireTime = createPayContext.Pay.CreateTime + int64(createPayContext.DaysUtilDue*86400)
			insert, err := dao.Payment.Ctx(ctx).Data(createPayContext.Pay).OmitEmpty().Insert(createPayContext.Pay)
			if err != nil {
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				return err
			}
			createPayContext.Pay.Id = id
			if createPayContext.Invoice.Id > 0 {
				_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
					dao.Invoice.Columns().PaymentId: createPayContext.Pay.PaymentId,
					dao.Invoice.Columns().GmtModify: gtime.Now(),
				}).Where(dao.Invoice.Columns().Id, createPayContext.Invoice.Id).OmitNil().Update()
				if err != nil {
					return err
				}
				invoice = query.GetInvoiceByInvoiceId(ctx, createPayContext.Invoice.InvoiceId)
			} else {
				invoice, err = handler.CreateProcessInvoiceForNewPayment(ctx, createPayContext.Invoice, createPayContext.Pay)
				if err != nil {
					return err
				}
			}
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

	gatewayInternalPayResult, err = api.GetGatewayServiceProvider(ctx, createPayContext.Pay.GatewayId).GatewayNewPayment(ctx, createPayContext)
	if err != nil {
		return nil, err
	}
	jsonData, err := gjson.Marshal(gatewayInternalPayResult)
	if err != nil {
		return nil, err
	}
	var automatic = 0
	if gatewayInternalPayResult.Status == consts.PaymentSuccess && createPayContext.PayImmediate {
		automatic = 1
	}
	createPayContext.Pay.PaymentData = string(jsonData)
	createPayContext.Pay.Status = int(gatewayInternalPayResult.Status)
	createPayContext.Pay.GatewayPaymentId = gatewayInternalPayResult.GatewayPaymentId
	createPayContext.Pay.GatewayPaymentIntentId = gatewayInternalPayResult.GatewayPaymentIntentId
	// unibee payment link
	paymentLink := link.GetPaymentLink(createPayContext.Pay.PaymentId)
	result, err := dao.Payment.Ctx(ctx).Data(g.Map{
		dao.Payment.Columns().PaymentData:            string(jsonData),
		dao.Payment.Columns().Automatic:              automatic,
		dao.Payment.Columns().Link:                   paymentLink,
		dao.Payment.Columns().GatewayLink:            gatewayInternalPayResult.Link,
		dao.Payment.Columns().GatewayPaymentId:       gatewayInternalPayResult.GatewayPaymentId,
		dao.Payment.Columns().GatewayPaymentIntentId: gatewayInternalPayResult.GatewayPaymentIntentId}).
		Where(dao.Payment.Columns().Id, createPayContext.Pay.Id).Where(dao.Payment.Columns().Status, consts.PaymentCreated).Update()
	if err != nil || result == nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return nil, err
	}
	gatewayInternalPayResult.Link = paymentLink
	createPayContext.Pay.Link = paymentLink
	gatewayInternalPayResult.Invoice = invoice
	gatewayInternalPayResult.Payment = createPayContext.Pay
	callback.GetPaymentCallbackServiceProvider(ctx, createPayContext.Pay.BizType).PaymentCreateCallback(ctx, createPayContext.Pay, gatewayInternalPayResult.Invoice)
	err = handler2.CreateOrUpdatePaymentTimelineForPayment(ctx, createPayContext.Pay, createPayContext.Pay.PaymentId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimelineForPayment error %s`, err.Error())
	}
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
		gatewayInternalPayResult.Invoice = query.GetInvoiceByInvoiceId(ctx, invoice.InvoiceId)
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

func CreateSubInvoicePaymentDefaultAutomatic(ctx context.Context, sub *entity.Subscription, invoice *entity.Invoice, gatewayId uint64, manualPayment bool, returnUrl string) (gatewayInternalPayResult *gateway_bean.GatewayNewPaymentResp, err error) {
	user := query.GetUserAccountById(ctx, sub.UserId)
	var email = ""
	if user != nil {
		email = user.Email
	}
	gateway := query.GetGatewayById(ctx, gatewayId)
	if gateway == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice gateway not found")
	}

	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	if merchant == nil {
		return nil, gerror.New("SubscriptionBillingCycleDunningInvoice merchantInfo not found")
	}
	invoice.Currency = strings.ToUpper(invoice.Currency)
	res, err := GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		PayImmediate: !manualPayment,
		CheckoutMode: manualPayment,
		Gateway:      gateway,
		Pay: &entity.Payment{
			SubscriptionId:    sub.SubscriptionId,
			ExternalPaymentId: sub.SubscriptionId,
			BizType:           consts.BizTypeSubscription,
			AuthorizeStatus:   consts.Authorized,
			UserId:            sub.UserId,
			GatewayId:         gateway.Id,
			TotalAmount:       invoice.TotalAmount,
			Currency:          strings.ToUpper(invoice.Currency),
			CryptoAmount:      invoice.CryptoAmount,
			CryptoCurrency:    invoice.CryptoCurrency,
			CountryCode:       sub.CountryCode,
			MerchantId:        sub.MerchantId,
			CompanyId:         merchant.CompanyId,
			Automatic:         1,
			ReturnUrl:         returnUrl,
			BillingReason:     invoice.InvoiceName,
			GasPayer:          sub.GasPayer,
		},
		ExternalUserId:       strconv.FormatUint(sub.UserId, 10),
		Email:                email,
		Invoice:              bean.SimplifyInvoice(invoice),
		Metadata:             map[string]interface{}{"BillingReason": invoice.InvoiceName},
		GatewayPaymentMethod: sub.GatewayDefaultPaymentMethod,
	})

	if err == nil && res.Payment != nil {
		if err == nil && res.Status != consts.PaymentSuccess {
			//need send invoice for authorised
			SendAuthorizedEmailBackground(sub, invoice, res.Payment)
		}
		if len(invoice.DiscountCode) > 0 {
			_, err = discount.UserDiscountApply(ctx, &discount.UserDiscountApplyReq{
				MerchantId:     sub.MerchantId,
				UserId:         sub.UserId,
				DiscountCode:   invoice.DiscountCode,
				SubscriptionId: sub.SubscriptionId,
				PaymentId:      res.Payment.PaymentId,
				InvoiceId:      invoice.InvoiceId,
				ApplyAmount:    invoice.DiscountAmount,
				Currency:       invoice.Currency,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return res, err
}

func HardDeletePayment(ctx context.Context, merchantId uint64, paymentId string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(paymentId) > 0, "invalid paymentId")
	one := query.GetPaymentByPaymentId(ctx, paymentId)
	if one != nil && len(one.InvoiceId) > 0 {
		_, err := dao.Invoice.Ctx(ctx).Where(dao.Invoice.Columns().InvoiceId, one.InvoiceId).Delete()
		if err != nil {
			return err
		}
	}
	_, err := dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().PaymentId, paymentId).Delete()
	return err
}

func SendAuthorizedEmailBackground(sub *entity.Subscription, invoice *entity.Invoice, payment *entity.Payment) {
	ctx := context.Background()
	go func() {
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
		oneUser := query.GetUserAccountById(ctx, sub.UserId)
		plan := query.GetPlanById(ctx, sub.PlanId)
		if plan != nil && oneUser != nil && merchant != nil {
			err = email2.SendTemplateEmail(ctx, merchant.Id, oneUser.Email, oneUser.TimeZone, email2.TemplateSubscriptionNeedAuthorized, "", &email2.TemplateVariable{
				UserName:            oneUser.FirstName + " " + oneUser.LastName,
				MerchantProductName: plan.PlanName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        query.GetMerchantCountryConfigName(ctx, payment.MerchantId, oneUser.CountryCode),
				PaymentAmount:       utility.ConvertCentToDollarStr(invoice.TotalAmount, invoice.Currency),
				Currency:            strings.ToUpper(invoice.Currency),
				PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail SendAuthorizedEmailBackground err:%s", err.Error())
			}
		}
	}()

}
