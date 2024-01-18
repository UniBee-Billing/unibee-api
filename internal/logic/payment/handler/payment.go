package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
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

type HandlePayReq struct {
	PaymentId      string
	ChannelPayId   string
	ChannelTradeNo string
	PayFee         int64
	PayStatusEnum  consts.PayStatusEnum
	PaidTime       *gtime.Time
	ReceiptFee     int64
	CaptureFee     int64
	Reason         string
	PaymentMethod  string
}

func HandlePayExpired(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", req)
	pay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}

	event.SaveTimeLine(ctx, entity.Timeline{
		BizType:   0,
		BizId:     pay.PaymentId,
		Fee:       pay.PaymentFee,
		EventType: event.Expird.Type,
		Event:     event.Expird.Desc,
		OpenApiId: pay.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s", pay.PaymentId, "Expird"),
	})

	return HandlePayFailure(ctx, &HandlePayReq{
		PaymentId:     req.PaymentId,
		PayStatusEnum: consts.PAY_FAILED,
		Reason:        "system cancel by expired",
	})
}

func HandleCaptureFailed(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", req)
	pay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}
	//交易事件记录
	event.SaveTimeLine(ctx, entity.Timeline{
		BizType:   0,
		BizId:     pay.PaymentId,
		Fee:       req.CaptureFee,
		EventType: event.CaptureFailed.Type,
		Event:     event.CaptureFailed.Desc,
		OpenApiId: pay.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.PaymentId, "CaptureFailed", req.ChannelPayId),
		Message:   req.Reason,
	})
	return nil
}

func HandlePayAuthorized(ctx context.Context, pay *entity.Payment) (err error) {
	g.Log().Infof(ctx, "HandlePayAuthorized, pay=%s", pay)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil")
		return errors.New("支付不存在")
	}
	if pay.AuthorizeStatus == consts.AUTHORIZED {
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayAuthorized, pay.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.AUTHORIZED, dao.Payment.Columns().ChannelPaymentId: pay.ChannelPaymentId},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID, dao.Payment.Columns().AuthorizeStatus: consts.WAITING_AUTHORIZED})
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
	g.Log().Infof(ctx, "HandlePayAuthorized sendResult err=%s", err)
	if err == nil {
		event.SaveTimeLine(ctx, entity.Timeline{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       pay.PaymentFee,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			OpenApiId: pay.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.ChannelPaymentId, "Authorised"),
		})
	}

	return err

}

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePayFailure, req=%s", req)
	pay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if pay == nil {
		g.Log().Infof(ctx, "pay null, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}
	if pay.Status == consts.PAY_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if pay.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	var refundFee int64 = 0
	if pay.AuthorizeStatus != consts.WAITING_AUTHORIZED {
		refundFee = pay.PaymentFee
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancelld, pay.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PAY_FAILED, dao.Payment.Columns().RefundFee: refundFee},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
	}

	g.Log().Infof(ctx, "HandlePayFailure sendResult err=%s", err)
	if err == nil {
		//交易事件记录
		event.SaveTimeLine(ctx, entity.Timeline{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       0,
			EventType: event.Cancelled.Type,
			Event:     event.Cancelled.Desc,
			OpenApiId: pay.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.PaymentId, "Cancelled"),
			Message:   req.Reason,
		})
	}
	return err
}

func HandlePaySuccess(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePaySuccess, req=%s", req)

	if req.PaidTime == nil {
		return errors.New("invalid param PaidTime is nil")
	}
	if len(req.PaymentId) == 0 {
		return errors.New("invalid param PaymentId is nil")
	}
	pay := query.GetPaymentByPaymentId(ctx, req.PaymentId)

	if pay == nil {
		g.Log().Infof(ctx, "pay not found, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if pay.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "merchantOrderNo:%s payment already success", req.PaymentId)
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaySuccess, pay.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().Status:                 consts.PAY_SUCCESS,
				dao.Payment.Columns().PaidTime:               req.PaidTime,
				dao.Payment.Columns().ChannelPaymentIntentId: req.ChannelPayId,
				dao.Payment.Columns().ChannelPaymentId:       req.ChannelTradeNo,
				dao.Payment.Columns().ReceiptFee:             req.ReceiptFee,
				dao.Payment.Columns().RefundFee:              pay.PaymentFee - req.ReceiptFee},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
	}

	g.Log().Infof(ctx, "HandlePaySuccess sendResult err=%s", err)

	if err == nil {
		//try {
		//交易事件记录
		event.SaveTimeLine(ctx, entity.Timeline{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       req.ReceiptFee,
			EventType: event.Settled.Type,
			Event:     event.Settled.Desc,
			OpenApiId: pay.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.PaymentId, "Settled"),
			Message:   req.Reason,
		})
		//} catch (Exception e) {
		//	e.printStackTrace();
		//	log.info("save_event exception:{}",e.toString());
		//}
		////            mqUtil.addToMq(MqTopicEnum.PaySuccess,pay.getId());
		//// 为提高支付结果反馈速度，同步处理业务订单支付成功
		//OverseaPay queryPay = overseaPayMapper.selectById(pay.getId());
		//if (queryPay == null) {
		//	pay.setPayStatus(PayStatusEnum.PAY_SUCCESS.getCode());
		//	pay.setPaidTime(Instant.ofEpochMilli(req.getPaidTime().getTime()).atZone(ZoneId.systemDefault()).toLocalDateTime());
		//	pay.setChannelPayId(req.getChannelPayId());
		//	pay.setChannelTradeNo(req.getChannelTradeNo());
		//	pay.setReceiptFee(req.getReceiptFee());
		//	queryPay = pay;
		//}
		//com.hk.utils.BusinessWrapper result = bizOrderPayCallbackProviderFactory.getBizOrderPayCallbackServiceProvider(queryPay.getBizType()).paySuccessCallback(queryPay,req.getPaidTime());
	}
	return err
}
