package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/logic/payment/event"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type HandleRefundReq struct {
	RefundId         string
	GatewayRefundId  string
	RefundAmount     int64
	RefundStatusEnum consts.RefundStatusEnum
	RefundTime       *gtime.Time
	Reason           string
}

func HandleRefundCancelled(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundFailure, req=%s", utility.MarshalToJsonString(req))
	if len(req.RefundId) == 0 {
		return gerror.New("invalid param refundNo")
	}
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.RefundId)
		return gerror.New("refund not found")
	}
	if one.Status == consts.RefundFailed {
		g.Log().Infof(ctx, "already failure")
		return nil
	}
	if one.Status == consts.RefundCancelled {
		g.Log().Infof(ctx, "already cancelled")
		return nil
	}
	if one.Status == consts.RefundSuccess {
		g.Log().Infof(ctx, "refund already success")
		return gerror.New("refund already success")
	}
	payment := query.GetPaymentByPaymentId(ctx, one.RefundId)
	if payment == nil {
		g.Log().Infof(ctx, "pay is nil, refundId=%s", one.RefundId)
		return gerror.New("payment not found")
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundFailed, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().Status: consts.RefundCancelled, dao.Refund.Columns().RefundComment: req.Reason},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.RefundCreated})
			if err != nil || result == nil {
				//_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				//_ = transaction.Rollback()
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
		return err
	} else {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundCancelCallback(ctx, payment, one)
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       one.RefundAmount,
			EventType: event.RefundFailed.Type,
			Event:     event.RefundFailed.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "RefundFailed", one.RefundId),
			Message:   req.Reason,
		})
		err = CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
	}

	return err
}

func HandleRefundFailure(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundFailure, req=%s", utility.MarshalToJsonString(req))
	if len(req.RefundId) == 0 {
		return gerror.New("invalid param refundNo")
	}
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.RefundId)
		return gerror.New("refund not found")
	}
	if one.Status == consts.RefundFailed {
		g.Log().Infof(ctx, "already failure")
		return nil
	}
	if one.Status == consts.RefundCancelled {
		g.Log().Infof(ctx, "already cancelled")
		return nil
	}
	if one.Status == consts.RefundSuccess {
		g.Log().Infof(ctx, "refund already success")
		return gerror.New("refund already success")
	}
	payment := query.GetPaymentByPaymentId(ctx, one.RefundId)
	if payment == nil {
		g.Log().Infof(ctx, "pay is nil, refundId=%s", one.RefundId)
		return gerror.New("payment not found")
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundFailed, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().Status: consts.RefundFailed, dao.Refund.Columns().RefundComment: req.Reason},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.RefundCreated})
			if err != nil || result == nil {
				//_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				//_ = transaction.Rollback()
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
		return err
	} else {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundFailureCallback(ctx, payment, one)
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       one.RefundAmount,
			EventType: event.RefundFailed.Type,
			Event:     event.RefundFailed.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "RefundFailed", one.RefundId),
			Message:   req.Reason,
		})
		err = CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
	}

	return err
}

func HandleRefundSuccess(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundSuccess, req=%s", utility.MarshalToJsonString(req))
	if len(req.RefundId) == 0 {
		return gerror.New("invalid param refundNo")
	}
	if len(req.RefundId) == 0 && req.RefundAmount > 0 {
		return gerror.New("invalid param RefundAmount, should > 0")
	}
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, refundId=%s", req.RefundId)
		return gerror.New("refund not found")
	}
	if one.Status == consts.RefundSuccess {
		g.Log().Infof(ctx, "refund already success")
		return nil
	}
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", one.PaymentId)
		return gerror.New("payment not found")
	}
	var refundAt = gtime.Now().Timestamp()
	if req.RefundTime != nil {
		refundAt = req.RefundTime.Timestamp()
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundSuccess, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().Status: consts.RefundSuccess, dao.Refund.Columns().RefundTime: refundAt},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.RefundCreated})
			if err != nil || result == nil {
				//_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				//_ = transaction.Rollback()
				return err
			}
			update, err := transaction.Update(dao.Payment.Table(), "refund_amount = refund_amount + ?", "id = ? AND ? >= 0 AND total_amount - refund_amount >= ?", one.RefundAmount, payment.Id, one.RefundAmount, one.RefundAmount)
			if err != nil || update == nil {
				//_ = transaction.Rollback()
				return err
			}
			payAffected, err := update.RowsAffected()
			g.Log().Printf(ctx, "HandleRefundSuccess Blank incrTotalRefundFee, updateCount=%v", payAffected)
			if err != nil || payAffected != 1 {
				//_ = transaction.Rollback()
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
		return err
	} else {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundSuccessCallback(ctx, payment, one)
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       one.RefundAmount,
			EventType: event.Refunded.Type,
			Event:     event.Refunded.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%d_%s_%s", payment.Status, "Refunded", one.RefundId),
			Message:   req.Reason,
		})
		err = CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
		}
	}

	return err
}

func HandleRefundReversed(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundReversed, req=%s", utility.MarshalToJsonString(req))
	if len(req.RefundId) == 0 {
		return gerror.New("invalid param refundNo")
	}
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.RefundId)
		return gerror.New("refund not found")
	}
	if one.Status != consts.RefundCreated {
		g.Log().Infof(ctx, "Refund is success or failure")
		return nil
	}
	payment := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", one.PaymentId)
		return gerror.New("payment not found")
	}
	callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentRefundReverseCallback(ctx, payment, one)
	// todo mark 此异常流有争议暂时什么都不做，只记录明细
	event.SaveEvent(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     payment.PaymentId,
		Fee:       one.RefundAmount,
		EventType: event.RefundedReversed.Type,
		Event:     event.RefundedReversed.Desc,
		OpenApiId: one.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "RefundedReversed", one.RefundId),
		Message:   req.Reason,
	})
	err = CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
	}

	return nil
}

func HandleRefundWebhookEvent(ctx context.Context, gatewayRefundRo *gateway_bean.GatewayPaymentRefundResp) error {
	utility.Assert(len(gatewayRefundRo.GatewayRefundId) > 0, "gatewayRefundId not found")
	one := query.GetRefundByGatewayRefundId(ctx, gatewayRefundRo.GatewayRefundId)
	utility.Assert(one != nil, "refund not found")
	if gatewayRefundRo.Status == consts.RefundSuccess {
		err := HandleRefundSuccess(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
			RefundAmount:     gatewayRefundRo.RefundAmount,
			RefundStatusEnum: gatewayRefundRo.Status,
			RefundTime:       gatewayRefundRo.RefundTime,
			Reason:           gatewayRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	} else if gatewayRefundRo.Status == consts.RefundFailed {
		err := HandleRefundFailure(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
			RefundAmount:     gatewayRefundRo.RefundAmount,
			RefundStatusEnum: gatewayRefundRo.Status,
			RefundTime:       gatewayRefundRo.RefundTime,
			Reason:           gatewayRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	} else if gatewayRefundRo.Status == consts.RefundCancelled {
		err := HandleRefundCancelled(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
			RefundAmount:     gatewayRefundRo.RefundAmount,
			RefundStatusEnum: gatewayRefundRo.Status,
			RefundTime:       gatewayRefundRo.RefundTime,
			Reason:           gatewayRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	} else if gatewayRefundRo.Status == consts.RefundReverse {
		err := HandleRefundReversed(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			GatewayRefundId:  gatewayRefundRo.GatewayRefundId,
			RefundAmount:     gatewayRefundRo.RefundAmount,
			RefundStatusEnum: gatewayRefundRo.Status,
			RefundTime:       gatewayRefundRo.RefundTime,
			Reason:           gatewayRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateOrUpdateRefundByDetail(ctx context.Context, payment *entity.Payment, details *gateway_bean.GatewayPaymentRefundResp, uniqueId string) error {
	utility.Assert(len(details.GatewayRefundId) > 0, "GatewayRefundId is null")
	utility.Assert(payment != nil, "payment is null")

	one := query.GetRefundByGatewayRefundId(ctx, details.GatewayRefundId)

	if one == nil {
		//创建
		one = &entity.Refund{
			CompanyId:            payment.CompanyId,
			MerchantId:           payment.MerchantId,
			UserId:               payment.UserId,
			OpenApiId:            payment.OpenApiId,
			GatewayId:            payment.GatewayId,
			CountryCode:          payment.CountryCode,
			Currency:             details.Currency,
			PaymentId:            payment.PaymentId,
			RefundId:             utility.CreateRefundId(),
			RefundAmount:         details.RefundAmount,
			RefundComment:        details.Reason,
			Status:               int(details.Status),
			RefundTime:           details.RefundTime.Timestamp(),
			GatewayRefundId:      details.GatewayRefundId,
			RefundCommentExplain: details.Reason,
			UniqueId:             uniqueId,
			SubscriptionId:       payment.SubscriptionId,
			CreateTime:           gtime.Now().Timestamp(),
		}

		result, err := dao.Refund.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateRefundByDetail record insert failure %s`, err.Error())
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = id
	} else {
		//更新

		_, err := dao.Refund.Ctx(ctx).Data(g.Map{
			dao.Refund.Columns().CompanyId:            payment.CompanyId,
			dao.Refund.Columns().MerchantId:           payment.MerchantId,
			dao.Refund.Columns().UserId:               payment.UserId,
			dao.Refund.Columns().OpenApiId:            payment.OpenApiId,
			dao.Refund.Columns().GatewayId:            payment.GatewayId,
			dao.Refund.Columns().CountryCode:          payment.CountryCode,
			dao.Refund.Columns().Currency:             details.Currency,
			dao.Refund.Columns().PaymentId:            payment.PaymentId,
			dao.Refund.Columns().RefundAmount:         details.RefundAmount,
			dao.Refund.Columns().RefundComment:        details.Reason,
			dao.Refund.Columns().Status:               details.Status,
			dao.Refund.Columns().RefundTime:           details.RefundTime.Timestamp(),
			dao.Refund.Columns().GatewayRefundId:      details.GatewayRefundId,
			dao.Refund.Columns().RefundCommentExplain: details.Reason,
			dao.Refund.Columns().UniqueId:             uniqueId,
			dao.Refund.Columns().SubscriptionId:       payment.SubscriptionId,
			dao.Refund.Columns().GmtModify:            gtime.Now(),
		}).Where(dao.Refund.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("CreateOrUpdateRefundByDetail err:%s", update)
		//}
	}

	return nil
}

func CreateOrUpdatePaymentTimelineFromRefund(ctx context.Context, refund *entity.Refund, uniqueId string) error {
	one := query.GetPaymentTimeLineByUniqueId(ctx, uniqueId)

	var status = 0
	if refund.Status == consts.RefundSuccess {
		status = 1
	} else if refund.Status == consts.RefundFailed {
		status = 2
	} else if refund.Status == consts.RefundReverse {
		status = 3
	}

	if one == nil {
		//创建
		one = &entity.PaymentTimeline{
			MerchantId:     refund.MerchantId,
			UserId:         refund.UserId,
			SubscriptionId: refund.SubscriptionId,
			//InvoiceId:      refund.InvoiceId,
			UniqueId:     uniqueId,
			Currency:     refund.Currency,
			TotalAmount:  refund.RefundAmount,
			GatewayId:    refund.GatewayId,
			Status:       status,
			TimelineType: 1,
			CreateTime:   gtime.Now().Timestamp(),
		}

		result, err := dao.PaymentTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdatePaymentTimelineFromRefund record insert failure %s`, err.Error())
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(id)
	} else {
		//更新
		_, err := dao.PaymentTimeline.Ctx(ctx).Data(g.Map{
			dao.PaymentTimeline.Columns().MerchantId:     refund.MerchantId,
			dao.PaymentTimeline.Columns().UserId:         refund.UserId,
			dao.PaymentTimeline.Columns().SubscriptionId: refund.SubscriptionId,
			//dao.PaymentTimeline.Columns().InvoiceId:      refund.InvoiceId,
			dao.PaymentTimeline.Columns().Currency:    refund.Currency,
			dao.PaymentTimeline.Columns().TotalAmount: refund.RefundAmount,
			dao.PaymentTimeline.Columns().GatewayId:   refund.GatewayId,
			//dao.PaymentTimeline.Columns().PaymentId:      payment.PaymentId,
			dao.PaymentTimeline.Columns().GmtModify:    gtime.Now(),
			dao.PaymentTimeline.Columns().Status:       status,
			dao.PaymentTimeline.Columns().TimelineType: 1,
		}).Where(dao.PaymentTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
	}
	return nil
}
