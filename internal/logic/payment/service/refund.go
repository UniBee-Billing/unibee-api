package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"golang.org/x/text/currency"
	"strconv"
	"strings"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/logic/payment/event"
	"unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type NewPaymentRefundInternalReq struct {
	PaymentId        string                 `json:"path" dc:"PaymentId" v:"required"`
	ExternalRefundId string                 `json:"externalRefundId" dc:"ExternalRefundId" v:"required"`
	RefundAmount     int64                  `json:"refundAmount" dc:"RefundAmount, Cent" v:"required"`
	Currency         string                 `json:"currency" dc:"Currency"  v:"required"`
	Reason           string                 `json:"reason" dc:"Reason"`
	Metadata         map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func GatewayPaymentRefundCreate(ctx context.Context, req *NewPaymentRefundInternalReq) (refund *entity.Refund, err error) {
	utility.Assert(len(req.PaymentId) > 0, "invalid paymentId")
	g.Log().Infof(ctx, "GatewayPaymentRefundCreate:%s", req.PaymentId)
	utility.Assert(len(req.ExternalRefundId) > 0, "invalid merchantRefundId")
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.TotalAmount > 0, "TotalAmount fee error")
	req.Currency = strings.ToUpper(req.Currency)
	utility.Assert(strings.Compare(payment.Currency, req.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(payment.Status == consts.PaymentSuccess, "payment not success")

	gateway := query.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

	utility.Assert(req.RefundAmount > 0, "refund value should > 0")
	utility.Assert(req.RefundAmount <= payment.TotalAmount-payment.RefundAmount, fmt.Sprintf("Already refunded %v%s", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(payment.Currency))), utility.ConvertCentToDollarStr(payment.RefundAmount, payment.Currency)))

	gatewayInfo := api.GetGatewayWebhookServiceProviderByGatewayName(ctx, gateway.GatewayName).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway information not found")
	if gatewayInfo.QueueForRefund {
		pendingRefunds := query.GetPendingGatewayTypeRefundsByPaymentId(ctx, req.PaymentId)
		utility.Assert(len(pendingRefunds) == 0, "Other refund still in process")
	}

	redisKey := fmt.Sprintf("createRefund-paymentId:%s-bizId:%s", payment.PaymentId, req.ExternalRefundId)
	isDuplicatedInvoke := false

	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}

	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		utility.Assert(false, "Submit Too Fast")
	}

	var gatewayExchange *gateway_bean.GatewayCurrencyExchange
	var exchangeRefundAmount int64 = 0
	var exchangeRefundCurrency string
	{
		//CurrencyExchange
		var PaymentMetadata = make(map[string]interface{})
		if len(payment.MetaData) > 0 {
			_ = utility.UnmarshalFromJsonString(payment.MetaData, &PaymentMetadata)
		}
		if s, ok := PaymentMetadata[gateway_bean.GatewayCurrencyExchangeKey]; ok {
			if value, ok2 := s.(string); ok2 {
				if len(value) > 0 {
					_ = utility.UnmarshalFromJsonString(value, &gatewayExchange)
				}
			}
		}
		if gatewayExchange != nil {
			req.Metadata[gateway_bean.GatewayCurrencyExchangeKey] = utility.MarshalToJsonString(gatewayExchange)
			//exchangeRefundAmount = int64(float64(req.RefundAmount) * gatewayExchange.ExchangeRate)
			exchangeRefundAmount = utility.ExchangeCurrencyConvert(req.RefundAmount, gatewayExchange.FromCurrency, gatewayExchange.ToCurrency, gatewayExchange.ExchangeRate)
			exchangeRefundCurrency = gatewayExchange.ToCurrency
		}
	}

	var (
		one *entity.Refund
	)
	err = dao.Refund.Ctx(ctx).Where(entity.Refund{
		PaymentId:        req.PaymentId,
		ExternalRefundId: req.ExternalRefundId,
		BizType:          payment.BizType,
	}).OmitEmpty().Scan(&one)
	utility.Assert(err == nil && one == nil, "Duplicate Submit")

	refundId := utility.CreateRefundId()
	req.Metadata["PaymentId"] = payment.PaymentId
	req.Metadata["RefundId"] = refundId
	req.Metadata["MerchantId"] = strconv.FormatUint(payment.MerchantId, 10)

	one = &entity.Refund{
		SubscriptionId:   payment.SubscriptionId,
		CompanyId:        payment.CompanyId,
		MerchantId:       payment.MerchantId,
		ExternalRefundId: req.ExternalRefundId,
		BizType:          payment.BizType,
		PaymentId:        payment.PaymentId,
		RefundId:         refundId,
		RefundAmount:     req.RefundAmount,
		Status:           consts.RefundCreated,
		GatewayId:        payment.GatewayId,
		AppId:            payment.AppId,
		Currency:         payment.Currency,
		CountryCode:      payment.CountryCode,
		RefundComment:    req.Reason,
		UserId:           payment.UserId,
		MetaData:         utility.MarshalToJsonString(req.Metadata),
	}

	//create Refund Invoice
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundCreated, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			one.InvoiceId = utility.CreateInvoiceId()
			one.UniqueId = one.RefundId
			one.CreateTime = gtime.Now().Timestamp()
			insert, err := dao.Refund.Ctx(ctx).Data(one).OmitEmpty().Insert()
			if err != nil {
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				return err
			}
			one.Id = id
			_, err = dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().PaymentId, payment.PaymentId).Increment(dao.Payment.Columns().RefundAmount, req.RefundAmount)
			if err != nil {
				return err
			}
			_, err = handler2.CreateProcessInvoiceForNewPaymentRefund(ctx, invoice_compute.CreateInvoiceSimplifyForRefund(ctx, payment, one), one)
			if err != nil {
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
	if err != nil {
		return nil, err
	}
	gatewayResult, err := api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayRefund(ctx, gateway, &gateway_bean.GatewayNewPaymentRefundReq{
		Payment:                 payment,
		Refund:                  one,
		Gateway:                 gateway,
		GatewayCurrencyExchange: gatewayExchange,
		ExchangeRefundCurrency:  exchangeRefundCurrency,
		ExchangeRefundAmount:    exchangeRefundAmount,
	})
	if err != nil {
		// todo mark record err to db
		return nil, err
	}

	one.GatewayRefundId = gatewayResult.GatewayRefundId
	one.Status = int(gatewayResult.Status)
	one.Type = gatewayResult.Type
	_, err = dao.Refund.Ctx(ctx).Data(g.Map{
		dao.Refund.Columns().GatewayRefundId:       gatewayResult.GatewayRefundId,
		dao.Refund.Columns().Type:                  gatewayResult.Type,
		dao.Refund.Columns().RefundGatewaySequence: gatewayResult.RefundSequence,
	}).Where(dao.Refund.Columns().Id, one.Id).Where(dao.Refund.Columns().Status, consts.RefundCreated).Update()
	if err != nil {
		return nil, err
	}
	// send the payment status checker mq
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismqcmd.TopicRefundChecker.Topic,
		Tag:        redismqcmd.TopicRefundChecker.Tag,
		Body:       one.RefundId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})

	if err != nil {
		return nil, err
	} else {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundCreateCallback(ctx, payment, one)
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.SentForRefund.Type,
			Event:     event.SentForRefund.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "SentForRefund", one.RefundId),
		})
		err = handler.CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
		if gatewayResult.Status == consts.RefundSuccess {
			err = handler.HandleRefundSuccess(ctx, &handler.HandleRefundReq{
				RefundId:         one.RefundId,
				GatewayRefundId:  gatewayResult.GatewayRefundId,
				RefundAmount:     req.RefundAmount,
				RefundStatusEnum: gatewayResult.Status,
				RefundTime:       gtime.Now(),
				Reason:           req.Reason,
			})
			if err != nil {
				return nil, err
			}
		} else if gatewayResult.Status == consts.RefundFailed {
			err = handler.HandleRefundFailure(ctx, &handler.HandleRefundReq{
				RefundId:         one.RefundId,
				GatewayRefundId:  gatewayResult.GatewayRefundId,
				RefundAmount:     gatewayResult.RefundAmount,
				RefundStatusEnum: gatewayResult.Status,
				RefundTime:       gatewayResult.RefundTime,
				Reason:           gatewayResult.Reason,
			})
		}
	}
	return one, nil
}

func MarkPaymentRefundCreate(ctx context.Context, req *NewPaymentRefundInternalReq) (refund *entity.Refund, err error) {
	utility.Assert(len(req.PaymentId) > 0, "invalid paymentId")
	g.Log().Infof(ctx, "MarkPaymentRefundCreate:%s", req.PaymentId)
	utility.Assert(len(req.ExternalRefundId) > 0, "invalid merchantRefundId")
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.TotalAmount > 0, "TotalAmount fee error")
	req.Currency = strings.ToUpper(req.Currency)
	utility.Assert(strings.Compare(payment.Currency, req.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(payment.Status == consts.PaymentSuccess, "payment not success")

	gateway := query.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	if gateway.GatewayType != consts.GatewayTypeCrypto && gateway.GatewayType != consts.GatewayTypeWireTransfer {
		return nil, gerror.New("mark refund only support crypto or wire transfer")
	}

	utility.Assert(req.RefundAmount > 0, "refund value should > 0")
	utility.Assert(req.RefundAmount <= payment.TotalAmount, "refund value should <= TotalAmount value")

	redisKey := fmt.Sprintf("createRefund-paymentId:%s-bizId:%s", payment.PaymentId, req.ExternalRefundId)
	isDuplicatedInvoke := false

	defer func() {
		if !isDuplicatedInvoke {
			utility.ReleaseLock(ctx, redisKey)
		}
	}()

	if !utility.TryLock(ctx, redisKey, 15) {
		isDuplicatedInvoke = true
		utility.Assert(false, "Submit Too Fast")
	}

	var (
		one *entity.Refund
	)
	err = dao.Refund.Ctx(ctx).Where(entity.Refund{
		PaymentId:        req.PaymentId,
		ExternalRefundId: req.ExternalRefundId,
		BizType:          payment.BizType,
	}).OmitEmpty().Scan(&one)
	utility.Assert(err == nil && one == nil, "Duplicate Submit")

	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	refundId := utility.CreateRefundId()
	req.Metadata["PaymentId"] = payment.PaymentId
	req.Metadata["RefundId"] = refundId
	req.Metadata["MerchantId"] = strconv.FormatUint(payment.MerchantId, 10)

	one = &entity.Refund{
		SubscriptionId:   payment.SubscriptionId,
		CompanyId:        payment.CompanyId,
		MerchantId:       payment.MerchantId,
		ExternalRefundId: req.ExternalRefundId,
		BizType:          payment.BizType,
		PaymentId:        payment.PaymentId,
		RefundId:         refundId,
		RefundAmount:     req.RefundAmount,
		Status:           consts.RefundCreated,
		GatewayId:        payment.GatewayId,
		Type:             consts.RefundTypeMarked,
		AppId:            payment.AppId,
		Currency:         payment.Currency,
		CountryCode:      payment.CountryCode,
		RefundComment:    req.Reason,
		UserId:           payment.UserId,
		MetaData:         utility.MarshalToJsonString(req.Metadata),
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundCreated, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			one.UniqueId = one.RefundId
			one.CreateTime = gtime.Now().Timestamp()
			insert, err := dao.Refund.Ctx(ctx).Data(one).OmitEmpty().Insert(one)
			if err != nil {
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				return err
			}
			one.Id = id

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
	} else {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundCreateCallback(ctx, payment, one)
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.SentForRefund.Type,
			Event:     event.SentForRefund.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "SentForRefund", one.RefundId),
		})
		err = handler.CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
		err = handler.HandleRefundSuccess(ctx, &handler.HandleRefundReq{
			RefundId:         one.RefundId,
			GatewayRefundId:  one.RefundId,
			RefundAmount:     one.RefundAmount,
			RefundStatusEnum: consts.RefundSuccess,
			RefundTime:       gtime.Now(),
			Reason:           one.RefundComment,
		})
		if err != nil {
			return nil, err
		}
	}
	return one, nil
}

func HardDeleteRefund(ctx context.Context, merchantId uint64, refundId string) error {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(len(refundId) > 0, "invalid refundId")
	_, err := dao.Refund.Ctx(ctx).Where(dao.Refund.Columns().RefundId, refundId).Delete()
	return err
}
