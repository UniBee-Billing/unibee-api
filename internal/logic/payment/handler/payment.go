package handler

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismqcmd "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
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
	pay := query.GetOverseaPayByMerchantOrderNo(ctx, req.MerchantOrderNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}
	//OverseaPayEvent overseaPayEvent = new OverseaPayEvent();
	//overseaPayEvent.setBizType(0);
	//overseaPayEvent.setBizId(pay.getId());
	//overseaPayEvent.setFee(pay.getPaymentFee());
	//overseaPayEvent.setEventType(TradeEventTypeEnum.Expird.getId());
	//overseaPayEvent.setEvent(TradeEventTypeEnum.Expird.getDesc());
	//overseaPayEvent.setUniqueNo(merchantOrderNo+"_Expird");
	//boolean save = iOverseaPayEventService.save(overseaPayEvent);
	//Assert.isTrue(save,"save event failure");

	return HandlePayFailure(ctx, &HandlePayReq{
		MerchantOrderNo: req.MerchantOrderNo,
		PayStatusEnum:   consts.PAY_FAILED,
		Reason:          "system cancel by expired",
	})
}

func HandleCaptureFailed(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", req)
	pay := query.GetOverseaPayByMerchantOrderNo(ctx, req.MerchantOrderNo)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}
	//交易事件记录
	//OverseaPayEvent overseaPayEvent = new OverseaPayEvent();
	//overseaPayEvent.setBizType(0);
	//overseaPayEvent.setBizId(pay.Id);
	//overseaPayEvent.setFee(req.CaptureFee);
	//overseaPayEvent.setEventType(TradeEventTypeEnum.CaptureFailed.getId());
	//overseaPayEvent.setEvent(TradeEventTypeEnum.CaptureFailed.getDesc());
	//overseaPayEvent.setUniqueNo(merchantOrderNo+"_CaptureFailed_"+req.ChannelPayId);
	//overseaPayEvent.setMessage(req.Reason);
	//boolean save = iOverseaPayEventService.save(overseaPayEvent);
	return nil
}

func HandlePayAuthorized(ctx context.Context, pay *entity.OverseaPay) (err error) {
	g.Log().Infof(ctx, "HandlePayAuthorized, pay=%s", pay)
	if pay == nil {
		g.Log().Infof(ctx, "pay is nil")
		return errors.New("支付不存在")
	}
	if pay.AuthorizeStatus == consts.AUTHORIZED {
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayAuthorized, pay.Id), func(messageToSend *redismq.Message) redismq.TransactionStatus {
		err = dao.OverseaPay.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update("oversea_pay", g.Map{"authorize_status": consts.AUTHORIZED, "channel_trade_no": pay.ChannelTradeNo},
				g.Map{"id": pay.Id, "pay_status": consts.TO_BE_PAID, "authorize_status": consts.WAITING_AUTHORIZED})
			if err != nil || result == nil {
				_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				_ = transaction.Rollback()
				return err
			}
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction
		} else {
			return redismq.RollbackTransaction
		}
	})
	g.Log().Infof(ctx, "HandlePayAuthorized sendResult err=%s", err)
	if err != nil {
		//OverseaPayEvent overseaPayEvent = new OverseaPayEvent();
		//overseaPayEvent.setBizType(0);
		//overseaPayEvent.setBizId(overseaPay.getId());
		//overseaPayEvent.setFee(overseaPay.getPaymentFee());
		//overseaPayEvent.setEventType(TradeEventTypeEnum.Authorised.getId());
		//overseaPayEvent.setEvent(TradeEventTypeEnum.Authorised.getDesc());
		//overseaPayEvent.setUniqueNo(overseaPay.getChannelTradeNo()+"_Authorised");
		//boolean save = iOverseaPayEventService.save(overseaPayEvent);
	}

	return err

}

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePayFailure, req=%s", req)
	pay := query.GetOverseaPayByMerchantOrderNo(ctx, req.MerchantOrderNo)
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
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancelld, pay.Id), func(messageToSend *redismq.Message) redismq.TransactionStatus {
		err = dao.OverseaPay.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update("oversea_pay", g.Map{"pay_status": consts.PAY_FAILED, "refund_fee": refundFee},
				g.Map{"id": pay.Id, "pay_status": consts.TO_BE_PAID})
			if err != nil || result == nil {
				_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				_ = transaction.Rollback()
				return err
			}
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction
		} else {
			return redismq.RollbackTransaction
		}
	})
	if err != nil {
		return err
	}

	g.Log().Infof(ctx, "HandlePayFailure sendResult err=%s", err)
	if err == nil {
		//try {
		//	//交易事件记录 todo mark
		//	OverseaPayEvent overseaPayEvent = new OverseaPayEvent();
		//	overseaPayEvent.setBizType(0);
		//	overseaPayEvent.setBizId(pay.getId());
		//	overseaPayEvent.setFee(0L);
		//	overseaPayEvent.setEventType(TradeEventTypeEnum.Cancelled.getId());
		//	overseaPayEvent.setEvent(TradeEventTypeEnum.Cancelled.getDesc());
		//	overseaPayEvent.setUniqueNo(pay.getMerchantOrderNo()+"_Cancelled");
		//	overseaPayEvent.setMessage(req.getReason());
		//	boolean save = iOverseaPayEventService.save(overseaPayEvent);
		//	Assert.isTrue(save,"save event failure");
		//} catch (Exception e) {
		//	e.printStackTrace();
		//	log.info("save_event exception:{}",e.toString());
		//}
	}
	return err
}

func HandlePaySuccess(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePaySuccess, req=%s", req)

	if req.PaidTime == nil {
		return errors.New("invalid param PaidTime is nil")
	}
	pay := query.GetOverseaPayByMerchantOrderNo(ctx, req.MerchantOrderNo)

	if pay == nil {
		g.Log().Infof(ctx, "pay null, merchantOrderNo=%s", req.MerchantOrderNo)
		return errors.New("支付不存在")
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if pay.PayStatus == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaySuccess, pay.Id), func(messageToSend *redismq.Message) redismq.TransactionStatus {
		err = dao.OverseaPay.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update("oversea_pay", g.Map{
				"pay_status":       consts.PAY_SUCCESS,
				"paid_time":        req.PaidTime,
				"channel_pay_id":   req.ChannelPayId,
				"channel_trade_no": req.ChannelTradeNo,
				"receipt_fee":      req.ReceiptFee,
				"refund_fee":       pay.PaymentFee - req.ReceiptFee},
				g.Map{"id": pay.Id, "pay_status": consts.TO_BE_PAID})
			if err != nil || result == nil {
				_ = transaction.Rollback()
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				_ = transaction.Rollback()
				return err
			}
			return nil
		})
		if err == nil {
			return redismq.CommitTransaction
		} else {
			return redismq.RollbackTransaction
		}
	})
	if err != nil {
		return err
	}

	g.Log().Infof(ctx, "HandlePaySuccess sendResult err=%s", err)

	if err == nil {
		//try {
		//	//交易事件记录
		//	OverseaPayEvent overseaPayEvent = new OverseaPayEvent();
		//	overseaPayEvent.setBizType(0);
		//	overseaPayEvent.setBizId(pay.getId());
		//	overseaPayEvent.setFee(req.getReceiptFee());
		//	overseaPayEvent.setEventType(TradeEventTypeEnum.Settled.getId());
		//	overseaPayEvent.setEvent(TradeEventTypeEnum.Settled.getDesc());
		//	overseaPayEvent.setUniqueNo(pay.getMerchantOrderNo()+"_Settled");
		//	overseaPayEvent.setMessage(req.getReason());
		//	boolean save = iOverseaPayEventService.save(overseaPayEvent);
		//	Assert.isTrue(save,"save event failure");
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
