package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/logic/payment/event"
	"unibee/internal/logic/payment/handler"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type NewPaymentRefundInternalReq struct {
	PaymentId        string            `json:"path" dc:"PaymentId" v:"required"`
	ExternalRefundId string            `json:"externalRefundId" dc:"ExternalRefundId" v:"required"`
	RefundAmount     int64             `json:"refundAmount" dc:"RefundAmount, Cent" v:"required"`
	Currency         string            `json:"currency" dc:"Currency"  v:"required"`
	Reason           string            `json:"reason" dc:"Reason"`
	Metadata         map[string]string `json:"metadata" dc:"Metadata，Map"`
}

func GatewayPaymentRefundCreate(ctx context.Context, req *NewPaymentRefundInternalReq) (refund *entity.Refund, err error) {
	utility.Assert(len(req.PaymentId) > 0, "invalid paymentId")
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
	utility.Assert(req.RefundAmount <= payment.TotalAmount-payment.RefundAmount, "no enough amount can refund")

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
		req.Metadata = make(map[string]string)
	}
	refundId := utility.CreateRefundId()
	req.Metadata["PaymentId"] = payment.PaymentId
	req.Metadata["RefundId"] = refundId
	req.Metadata["MerchantId"] = strconv.FormatUint(payment.MerchantId, 10)

	one = &entity.Refund{
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
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundCreated, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
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

			_, err = dao.Payment.Ctx(ctx).Where(dao.Payment.Columns().PaymentId, payment.PaymentId).Increment(dao.Payment.Columns().RefundAmount, refund.RefundAmount)
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
	gatewayResult, err := api.GetGatewayServiceProvider(ctx, payment.GatewayId).GatewayRefund(ctx, payment, one)
	if err != nil {
		return nil, err
	}

	one.GatewayRefundId = gatewayResult.GatewayRefundId
	result, err := dao.Refund.Ctx(ctx).Data(g.Map{
		dao.Refund.Columns().GatewayRefundId: gatewayResult.GatewayRefundId}).
		Where(dao.Refund.Columns().Id, one.Id).Where(dao.Refund.Columns().Status, consts.RefundCreated).Update()
	if err != nil || result == nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return nil, err
	}

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
		}
	}
	return one, nil
}

func MarkPaymentRefundCreate(ctx context.Context, req *NewPaymentRefundInternalReq) (refund *entity.Refund, err error) {
	utility.Assert(len(req.PaymentId) > 0, "invalid paymentId")
	utility.Assert(len(req.ExternalRefundId) > 0, "invalid merchantRefundId")
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.TotalAmount > 0, "TotalAmount fee error")
	req.Currency = strings.ToUpper(req.Currency)
	utility.Assert(strings.Compare(payment.Currency, req.Currency) == 0, "refund currency not match the payment error")
	utility.Assert(payment.Status == consts.PaymentSuccess, "payment not success")

	gateway := query.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	if gateway.GatewayType != consts.GatewayTypeCrypto {
		return nil, gerror.New("mark refund only support crypto")
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
		req.Metadata = make(map[string]string)
	}
	refundId := utility.CreateRefundId()
	req.Metadata["PaymentId"] = payment.PaymentId
	req.Metadata["RefundId"] = refundId
	req.Metadata["MerchantId"] = strconv.FormatUint(payment.MerchantId, 10)

	one = &entity.Refund{
		CompanyId:        payment.CompanyId,
		MerchantId:       payment.MerchantId,
		ExternalRefundId: req.ExternalRefundId,
		BizType:          payment.BizType,
		PaymentId:        payment.PaymentId,
		RefundId:         refundId,
		RefundAmount:     req.RefundAmount,
		Status:           consts.RefundCreated,
		GatewayId:        payment.GatewayId,
		Type:             2, //mark refund type
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
				//_ = transaction.Rollback()
				return err
			}
			id, err := insert.LastInsertId()
			if err != nil {
				//_ = transaction.Rollback()
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
