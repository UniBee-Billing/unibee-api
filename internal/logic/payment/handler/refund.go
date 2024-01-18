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
	"go-oversea-pay/internal/logic/payment/event"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
)

type HandleRefundReq struct {
	MerchantRefundNo string
	ChannelRefundNo  string
	RefundFee        int64
	RefundStatusEnum consts.RefundStatusEnum
	RefundTime       *gtime.Time
	Reason           string
}

func HandleRefundFailure(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundFailure, req=%s", req)
	if len(req.MerchantRefundNo) == 0 {
		return gerror.New("invalid param refundNo")
	}
	one := query.GetRefundByMerchantRefundNo(ctx, req.MerchantRefundNo)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.MerchantRefundNo)
		return gerror.New("退款记录不存在")
	}
	if one.RefundStatus == consts.REFUND_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}
	if one.RefundStatus == consts.REFUND_SUCCESS {
		g.Log().Infof(ctx, "refund already success")
		return gerror.New("refund already success")
	}
	pay := query.GetPaymentByMerchantOrderNo(ctx, one.OutTradeNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", one.OutTradeNo)
		return gerror.New("支付记录不存在")
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundFailed, one.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().RefundStatus: consts.REFUND_FAILED, dao.Refund.Columns().RefundComment: req.Reason},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().RefundStatus: consts.REFUND_ING})
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
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     pay.Id,
			Fee:       one.RefundFee,
			EventType: event.RefundFailed.Type,
			Event:     event.RefundFailed.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.MerchantOrderNo, "RefundFailed", one.OutRefundNo),
			Message:   req.Reason,
		})
	}

	return err
}

func HandleRefundSuccess(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundSuccess, req=%s", req)
	if len(req.MerchantRefundNo) == 0 {
		return gerror.New("invalid param refundNo")
	}
	if len(req.MerchantRefundNo) == 0 && req.RefundFee > 0 {
		return gerror.New("invalid param RefundFee, should > 0")
	}
	one := query.GetRefundByMerchantRefundNo(ctx, req.MerchantRefundNo)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.MerchantRefundNo)
		return gerror.New("退款记录不存在")
	}
	if one.RefundStatus == consts.REFUND_SUCCESS {
		g.Log().Infof(ctx, "refund already success")
		return nil
	}
	pay := query.GetPaymentByMerchantOrderNo(ctx, one.OutTradeNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", one.OutTradeNo)
		return gerror.New("支付记录不存在")
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
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicRefundSuccess, one.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Refund.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Refund.Table(), g.Map{dao.Refund.Columns().RefundStatus: consts.REFUND_SUCCESS, dao.Refund.Columns().RefundTime: req.RefundTime},
				g.Map{dao.Refund.Columns().Id: one.Id, dao.Refund.Columns().RefundStatus: consts.REFUND_ING})
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
			update, err := transaction.Update(dao.Payment.Table(), "refund_fee = refund_fee + ?", "id = ? AND ? >= 0 AND payment_fee - refund_fee >= ?", one.RefundFee, pay.Id, one.RefundFee, one.RefundFee)
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
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     pay.Id,
			Fee:       one.RefundFee,
			EventType: event.Refunded.Type,
			Event:     event.Refunded.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.MerchantOrderNo, "Refunded", one.OutRefundNo),
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
	}

	return err
}

func HandleRefundReversed(ctx context.Context, req *HandleRefundReq) (err error) {
	g.Log().Infof(ctx, "HandleRefundReversed, req=%s", req)
	if len(req.MerchantRefundNo) == 0 {
		return gerror.New("invalid param refundNo")
	}
	one := query.GetRefundByMerchantRefundNo(ctx, req.MerchantRefundNo)
	if one == nil {
		g.Log().Infof(ctx, "refund is nil, merchantOrderNo=%s", req.MerchantRefundNo)
		return gerror.New("退款记录不存在")
	}
	if one.RefundStatus == consts.REFUND_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}
	pay := query.GetPaymentByMerchantOrderNo(ctx, one.OutTradeNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", one.OutTradeNo)
		return gerror.New("支付记录不存在")
	}
	// todo mark 此异常流有争议暂时什么都不做，只记录明细
	event.SaveEvent(ctx, entity.OverseaPayEvent{
		BizType:   0,
		BizId:     pay.Id,
		Fee:       one.RefundFee,
		EventType: event.RefundedReversed.Type,
		Event:     event.RefundedReversed.Desc,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.MerchantOrderNo, "RefundedReversed", one.OutRefundNo),
		Message:   req.Reason,
	})

	return nil
}
