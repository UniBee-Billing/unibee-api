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
	MerchantOrderNo string
	ChannelPayId    string
	ChannelTradeNo  string
	PayFee          int64
	PayStatusEnum   consts.PayStatusEnum
	PaidTime        *gtime.Time
	ReceiptFee      int64
	CaptureFee      int64
	Reason          string
	PaymentMethod   string
}

func HandlePayExpired(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", req)
	pay := query.GetPaymentByMerchantOrderNo(ctx, req.MerchantOrderNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}

	event.SaveEvent(ctx, entity.OverseaPayEvent{
		BizType:   0,
		BizId:     pay.Id,
		Fee:       pay.PaymentFee,
		EventType: event.Expird.Type,
		Event:     event.Expird.Desc,
		UniqueNo:  fmt.Sprintf("%s_%s", pay.MerchantOrderNo, "Expird"),
	})

	return HandlePayFailure(ctx, &HandlePayReq{
		MerchantOrderNo: req.MerchantOrderNo,
		PayStatusEnum:   consts.PAY_FAILED,
		Reason:          "system cancel by expired",
	})
}

func HandleCaptureFailed(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", req)
	pay := query.GetPaymentByMerchantOrderNo(ctx, req.MerchantOrderNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}
	//交易事件记录
	event.SaveEvent(ctx, entity.OverseaPayEvent{
		BizType:   0,
		BizId:     pay.Id,
		Fee:       req.CaptureFee,
		EventType: event.CaptureFailed.Type,
		Event:     event.CaptureFailed.Desc,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", pay.MerchantOrderNo, "CaptureFailed", req.ChannelPayId),
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
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.AUTHORIZED, dao.Payment.Columns().ChannelTradeNo: pay.ChannelTradeNo},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().PayStatus: consts.TO_BE_PAID, dao.Payment.Columns().AuthorizeStatus: consts.WAITING_AUTHORIZED})
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
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     pay.Id,
			Fee:       pay.PaymentFee,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.ChannelTradeNo, "Authorised"),
		})
	}

	return err

}

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePayFailure, req=%s", req)
	pay := query.GetPaymentByMerchantOrderNo(ctx, req.MerchantOrderNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay null, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}
	if pay.PayStatus == consts.PAY_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if pay.PayStatus == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	var refundFee int64 = 0
	if pay.AuthorizeStatus != consts.WAITING_AUTHORIZED {
		refundFee = pay.PaymentFee
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancelld, pay.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().PayStatus: consts.PAY_FAILED, dao.Payment.Columns().RefundFee: refundFee},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().PayStatus: consts.TO_BE_PAID})
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
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     pay.Id,
			Fee:       0,
			EventType: event.Cancelled.Type,
			Event:     event.Cancelled.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.MerchantOrderNo, "Cancelled"),
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
	if len(req.MerchantOrderNo) == 0 {
		return errors.New("invalid param MerchantOrderNo is nil")
	}
	pay := query.GetPaymentByMerchantOrderNo(ctx, req.MerchantOrderNo)

	if pay == nil {
		g.Log().Infof(ctx, "pay not found, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if pay.PayStatus == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "merchantOrderNo:%s payment already success", req.MerchantOrderNo)
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaySuccess, pay.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().PayStatus:      consts.PAY_SUCCESS,
				dao.Payment.Columns().PaidTime:       req.PaidTime,
				dao.Payment.Columns().ChannelPayId:   req.ChannelPayId,
				dao.Payment.Columns().ChannelTradeNo: req.ChannelTradeNo,
				dao.Payment.Columns().ReceiptFee:     req.ReceiptFee,
				dao.Payment.Columns().RefundFee:      pay.PaymentFee - req.ReceiptFee},
				g.Map{dao.Payment.Columns().Id: pay.Id, dao.Payment.Columns().PayStatus: consts.TO_BE_PAID})
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
		event.SaveEvent(ctx, entity.OverseaPayEvent{
			BizType:   0,
			BizId:     pay.Id,
			Fee:       req.ReceiptFee,
			EventType: event.Settled.Type,
			Event:     event.Settled.Desc,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.MerchantOrderNo, "Settled"),
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
