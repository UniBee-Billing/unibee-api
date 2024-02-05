package handler

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismqcmd "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/payment/event"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

type HandleRefundReq struct {
	RefundId         string
	ChannelRefundId  string
	RefundFee        int64
	RefundStatusEnum consts.RefundStatusEnum
	RefundTime       *gtime.Time
	Reason           string
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
	if one.Status == consts.REFUND_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}
	if one.Status == consts.REFUND_SUCCESS {
		g.Log().Infof(ctx, "refund already success")
		return gerror.New("refund already success")
	}
	pay := query.GetPaymentByPaymentId(ctx, one.RefundId)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, refundId=%s", one.RefundId)
		return gerror.New("payment not found")
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundFailed, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().Status: consts.REFUND_FAILED, dao.Refund.Columns().RefundComment: req.Reason},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.REFUND_ING})
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
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       one.RefundAmount,
			EventType: event.RefundFailed.Type,
			Event:     event.RefundFailed.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.PaymentId, "RefundFailed", one.RefundId),
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
	if len(req.RefundId) == 0 && req.RefundFee > 0 {
		return gerror.New("invalid param RefundAmount, should > 0")
	}
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, refundId=%s", req.RefundId)
		return gerror.New("refund not found")
	}
	if one.Status == consts.REFUND_SUCCESS {
		g.Log().Infof(ctx, "refund already success")
		return nil
	}
	pay := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", one.PaymentId)
		return gerror.New("payment not found")
	}
	//if (refund.getRefundComment().equals("手动触发重复支付单退款")) {
	//	int updateCount = overseaRefundMapper.update(new OverseaRefund(),
	//		new UpdateWrapper<OverseaRefund>().lambda()
	//	.set(OverseaRefund::getRefundStatus, RefundStatusEnum.REFUND_SUCCESS.getCode())
	//	.set(OverseaRefund::getRefundTime, req.getRefundTime())
	//	.eq(OverseaRefund::getOutRefundNo, outRefundNo)
	//	.eq(OverseaRefund::getRefundStatus, RefundStatusEnum.REFUND_ING.getCode())
	//);
	//	log.info("update refund status to REFUND_SUCCESS, updateCount={}", updateCount);
	//	if (updateCount != 1) {
	//		return BusinessWrapper.failWithMessage("手动触发重复支付单退款失败");
	//	}
	//	return BusinessWrapper.success();
	//}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundSuccess, one.RefundId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().Status: consts.REFUND_SUCCESS, dao.Refund.Columns().RefundTime: req.RefundTime},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().Status: consts.REFUND_ING})
			if err != nil || result == nil {
				//_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				//_ = transaction.Rollback()
				return err
			}
			//支付单补充退款金额
			update, err := transaction.Update(dao.Payment.Table(), "refund_amount = refund_amount + ?", "id = ? AND ? >= 0 AND total_amount - refund_amount >= ?", one.RefundAmount, pay.Id, one.RefundAmount, one.RefundAmount)
			if err != nil || update == nil {
				//_ = transaction.Rollback()
				return err
			}
			payAffected, err := update.RowsAffected()
			g.Log().Printf(ctx, "HandleRefundSuccess Blank incrTotalRefundFee, updateCount=%s", payAffected)
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
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       one.RefundAmount,
			EventType: event.Refunded.Type,
			Event:     event.Refunded.Desc,
			OpenApiId: one.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.Status, "Refunded", one.RefundId),
			Message:   req.Reason,
		})

		//producerWrapper.send(new Message(MqTopicEnum.RefundSuccess,refund.getId()));
		////            mqUtil.addToMq(MqTopicEnum.RefundSuccess,refund.getId());
		//// 为提高退款结果反馈速度，同步处理业务订单支付成功
		//OverseaRefund queryRefund = overseaRefundMapper.selectById(refund.getId());
		//if (queryRefund == null) {
		//	refund.setRefundStatus(RefundStatusEnum.REFUND_SUCCESS.getCode());
		//	refund.setRefundTime(Instant.ofEpochMilli(req.getRefundTime().getTime()).atZone(ZoneId.systemDefault()).toLocalDateTime());
		//	queryRefund = refund;
		//}
		//com.hk.utils.BusinessWrapper result = bizOrderPayCallbackProviderFactory.getBizOrderPayCallbackServiceProvider(queryRefund.getBizType()).refundSuccessCallback(queryRefund,req.getRefundTime());
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
	if one.Status != consts.REFUND_ING {
		g.Log().Infof(ctx, "Refund is success or failure")
		return nil
	}
	pay := query.GetPaymentByPaymentId(ctx, one.PaymentId)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", one.PaymentId)
		return gerror.New("payment not found")
	}
	// todo mark 此异常流有争议暂时什么都不做，只记录明细
	event.SaveEvent(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     pay.PaymentId,
		Fee:       one.RefundAmount,
		EventType: event.RefundedReversed.Type,
		Event:     event.RefundedReversed.Desc,
		OpenApiId: one.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.PaymentId, "RefundedReversed", one.RefundId),
		Message:   req.Reason,
	})
	err = CreateOrUpdatePaymentTimelineFromRefund(ctx, one, one.RefundId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimelineFromRefund error %s`, err.Error())
	}

	return nil
}

func HandleRefundWebhookEvent(ctx context.Context, channelRefundRo *ro.OutPayRefundRo) error {
	utility.Assert(len(channelRefundRo.ChannelRefundId) > 0, "channelRefundId not found")
	one := query.GetRefundByChannelRefundId(ctx, channelRefundRo.ChannelRefundId)
	utility.Assert(one != nil, "refund not found")
	if channelRefundRo.Status == consts.REFUND_SUCCESS {
		err := HandleRefundSuccess(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			ChannelRefundId:  channelRefundRo.ChannelRefundId,
			RefundFee:        channelRefundRo.RefundFee,
			RefundStatusEnum: channelRefundRo.Status,
			RefundTime:       channelRefundRo.RefundTime,
			Reason:           channelRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	} else if channelRefundRo.Status == consts.REFUND_FAILED {
		err := HandleRefundFailure(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			ChannelRefundId:  channelRefundRo.ChannelRefundId,
			RefundFee:        channelRefundRo.RefundFee,
			RefundStatusEnum: channelRefundRo.Status,
			RefundTime:       channelRefundRo.RefundTime,
			Reason:           channelRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	} else if channelRefundRo.Status == consts.REFUND_REVERSE {
		err := HandleRefundReversed(ctx, &HandleRefundReq{
			RefundId:         one.RefundId,
			ChannelRefundId:  channelRefundRo.ChannelRefundId,
			RefundFee:        channelRefundRo.RefundFee,
			RefundStatusEnum: channelRefundRo.Status,
			RefundTime:       channelRefundRo.RefundTime,
			Reason:           channelRefundRo.Reason,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateOrUpdateRefundByDetail(ctx context.Context, payment *entity.Payment, details *ro.OutPayRefundRo, uniqueId string) error {
	utility.Assert(len(details.ChannelPaymentId) > 0, "paymentId is null")
	utility.Assert(payment != nil, "payment is null")

	one := query.GetRefundByChannelRefundId(ctx, details.ChannelRefundId)

	if one == nil {
		//创建
		one = &entity.Refund{
			CompanyId:            payment.CompanyId,
			MerchantId:           payment.MerchantId,
			UserId:               payment.UserId,
			OpenApiId:            payment.OpenApiId,
			ChannelId:            payment.ChannelId,
			CountryCode:          payment.CountryCode,
			Currency:             details.Currency,
			PaymentId:            payment.PaymentId,
			RefundId:             utility.CreateRefundId(),
			RefundAmount:         details.RefundFee,
			RefundComment:        details.Reason,
			Status:               int(details.Status),
			RefundTime:           details.RefundTime,
			ChannelRefundId:      details.ChannelRefundId,
			RefundCommentExplain: details.Reason,
			UniqueId:             uniqueId,
			SubscriptionId:       payment.SubscriptionId,
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
			dao.Refund.Columns().ChannelId:            payment.ChannelId,
			dao.Refund.Columns().CountryCode:          payment.CountryCode,
			dao.Refund.Columns().Currency:             details.Currency,
			dao.Refund.Columns().PaymentId:            payment.PaymentId,
			dao.Refund.Columns().RefundAmount:         details.RefundFee,
			dao.Refund.Columns().RefundComment:        details.Reason,
			dao.Refund.Columns().Status:               details.Status,
			dao.Refund.Columns().RefundTime:           details.RefundTime,
			dao.Refund.Columns().ChannelRefundId:      details.ChannelRefundId,
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
	if refund.Status == consts.REFUND_SUCCESS {
		status = 1
	} else if refund.Status == consts.REFUND_FAILED {
		status = 2
	} else if refund.Status == consts.REFUND_REVERSE {
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
			ChannelId:    refund.ChannelId,
			Status:       status,
			TimelineType: 1,
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
			dao.PaymentTimeline.Columns().ChannelId:   refund.ChannelId,
			//dao.PaymentTimeline.Columns().PaymentId:      payment.PaymentId,
			dao.PaymentTimeline.Columns().GmtModify:    gtime.Now(),
			dao.PaymentTimeline.Columns().Status:       status,
			dao.PaymentTimeline.Columns().TimelineType: 1,
		}).Where(dao.PaymentTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("CreateOrUpdatePaymentTimelineFromRefund err:%s", update)
		//}
	}
	return nil
}
