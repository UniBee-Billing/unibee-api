package handler

import (
	"context"
	"errors"
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
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
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

	event.SaveTimeLine(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     pay.PaymentId,
		Fee:       pay.PaymentFee,
		EventType: event.Expird.Type,
		Event:     event.Expird.Desc,
		OpenApiId: pay.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s", pay.PaymentId, "Expird"),
	})

	err = CreateOrUpdatePaymentTimeline(ctx, pay, pay.PaymentId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
	}

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
	event.SaveTimeLine(ctx, entity.PaymentEvent{
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
		event.SaveTimeLine(ctx, entity.PaymentEvent{
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
		event.SaveTimeLine(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     pay.PaymentId,
			Fee:       0,
			EventType: event.Cancelled.Type,
			Event:     event.Cancelled.Desc,
			OpenApiId: pay.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", pay.PaymentId, "Cancelled"),
			Message:   req.Reason,
		})
		err := CreateOrUpdatePaymentTimeline(ctx, pay, pay.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
		}
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
		event.SaveTimeLine(ctx, entity.PaymentEvent{
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
		err := CreateOrUpdatePaymentTimeline(ctx, pay, pay.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
		}
	}
	return err
}

func HandlePaymentWebhookEvent(ctx context.Context, eventType string, details *ro.OutPayRo) error {
	_, err := CreateOrUpdatePaymentByDetail(ctx, details, details.ChannelPaymentId)
	if err != nil {
		return err
	}

	if details.ChannelInvoiceDetail != nil && details.Status == consts.PAY_SUCCESS && details.ChannelSubscriptionDetail != nil {
		err = handler.HandleSubscriptionPaymentSuccess(ctx, &handler.SubscriptionPaymentSuccessWebHookReq{
			ChannelPaymentId:      details.ChannelPaymentId,
			ChannelSubscriptionId: details.ChannelInvoiceDetail.ChannelSubscriptionId,
			ChannelInvoiceId:      details.ChannelInvoiceDetail.ChannelInvoiceId,
			ChannelUpdateId:       details.ChannelUpdateId,
			Status:                details.ChannelSubscriptionDetail.Status,
			ChannelStatus:         details.ChannelInvoiceDetail.ChannelStatus,
			Data:                  details.ChannelSubscriptionDetail.Data,
			ChannelItemData:       details.ChannelSubscriptionDetail.ChannelItemData,
			CancelAtPeriodEnd:     details.ChannelSubscriptionDetail.CancelAtPeriodEnd,
			CurrentPeriodEnd:      details.ChannelSubscriptionDetail.CurrentPeriodEnd,
			CurrentPeriodStart:    details.ChannelSubscriptionDetail.CurrentPeriodStart,
			TrialEnd:              details.ChannelSubscriptionDetail.TrialEnd,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateOrUpdatePaymentByDetail(ctx context.Context, details *ro.OutPayRo, uniqueId string) (*entity.Payment, error) {
	utility.Assert(len(details.ChannelInvoiceId) > 0, "invoice id is null")
	var subscriptionId string
	var merchantId int64
	var channelId int64
	var userId int64
	var invoiceId string
	var countryCode string
	if len(details.ChannelInvoiceId) > 0 {
		invoice := query.GetInvoiceByChannelInvoiceId(ctx, details.ChannelInvoiceId)
		if invoice != nil {
			invoiceId = invoice.InvoiceId
			subscriptionId = invoice.SubscriptionId
			if len(subscriptionId) > 0 {
				sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
				countryCode = sub.CountryCode
			}
			merchantId = invoice.MerchantId
			channelId = invoice.ChannelId
			userId = invoice.UserId
		}
	}
	if details.ChannelSubscriptionDetail != nil {
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, details.ChannelSubscriptionDetail.ChannelSubscriptionId)
		subscriptionId = sub.SubscriptionId
		countryCode = sub.CountryCode
	}
	one := query.GetPaymentByChannelPaymentId(ctx, details.ChannelPaymentId)

	if one == nil {
		//创建
		one = &entity.Payment{
			MerchantId:             merchantId,
			UserId:                 userId,
			CountryCode:            countryCode,
			PaymentId:              utility.CreatePaymentId(),
			Currency:               details.Currency,
			PaymentFee:             details.PayFee,
			RefundFee:              details.TotalRefundFee,
			ReceiptFee:             details.PayFee,
			Status:                 details.Status,
			AuthorizeStatus:        details.CaptureStatus,
			ChannelId:              channelId,
			ChannelPaymentFee:      details.PayFee,
			ChannelPaymentIntentId: details.ChannelPaymentId,
			ChannelPaymentId:       details.ChannelPaymentId,
			CreateTime:             details.CreateTime,
			CancelTime:             details.CancelTime,
			PaidTime:               details.PayTime,
			ChannelInvoiceId:       details.ChannelInvoiceId,
			PaymentData:            details.CancelReason,
			UniqueId:               uniqueId,
			SubscriptionId:         subscriptionId,
			InvoiceId:              invoiceId,
		}
		result, err := dao.Payment.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdatePaymentByDetail record insert failure %s`, err.Error())
			return nil, err
		}
		id, _ := result.LastInsertId()
		one.Id = id
	} else {
		//更新
		update, err := dao.Payment.Ctx(ctx).Data(g.Map{
			dao.Payment.Columns().MerchantId:             merchantId,
			dao.Payment.Columns().MerchantId:             merchantId,
			dao.Payment.Columns().UserId:                 userId,
			dao.Payment.Columns().CountryCode:            countryCode,
			dao.Payment.Columns().PaymentId:              utility.CreatePaymentId(),
			dao.Payment.Columns().Currency:               details.Currency,
			dao.Payment.Columns().PaymentFee:             details.PayFee,
			dao.Payment.Columns().RefundFee:              details.TotalRefundFee,
			dao.Payment.Columns().ReceiptFee:             details.PayFee,
			dao.Payment.Columns().Status:                 details.Status,
			dao.Payment.Columns().AuthorizeStatus:        details.CaptureStatus,
			dao.Payment.Columns().ChannelId:              channelId,
			dao.Payment.Columns().ChannelPaymentFee:      details.PayFee,
			dao.Payment.Columns().ChannelPaymentIntentId: details.ChannelPaymentId,
			dao.Payment.Columns().ChannelPaymentId:       details.ChannelPaymentId,
			dao.Payment.Columns().CreateTime:             details.CreateTime,
			dao.Payment.Columns().CancelTime:             details.CancelTime,
			dao.Payment.Columns().PaidTime:               details.PayTime,
			dao.Payment.Columns().ChannelInvoiceId:       details.ChannelInvoiceId,
			dao.Payment.Columns().PaymentData:            details.CancelReason,
			dao.Payment.Columns().UniqueId:               uniqueId,
			dao.Payment.Columns().SubscriptionId:         subscriptionId,
			dao.Payment.Columns().InvoiceId:              invoiceId,
			dao.Invoice.Columns().GmtModify:              gtime.Now(),
		}).Where(dao.Payment.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return nil, err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return nil, gerror.Newf("CreateOrUpdatePaymentByDetail err:%s", update)
		}
	}
	err := CreateOrUpdatePaymentTimeline(ctx, one, one.PaymentId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
	}
	return one, nil
}

func CreateOrUpdatePaymentTimeline(ctx context.Context, payment *entity.Payment, uniqueId string) error {
	one := query.GetPaymentTimeLineByUniqueId(ctx, uniqueId)

	var status = 0
	if payment.Status == consts.PAY_SUCCESS {
		status = 1
	} else if payment.Status == consts.PAY_FAILED {
		status = 2
	}
	var timeLineType = 0
	if payment.PaymentFee > 0 {
		timeLineType = 0
	} else if payment.PaymentFee < 0 {
		timeLineType = 1
	}
	if one == nil {
		//创建
		one = &entity.PaymentTimeline{
			MerchantId:     payment.MerchantId,
			UserId:         payment.UserId,
			SubscriptionId: payment.SubscriptionId,
			InvoiceId:      payment.InvoiceId,
			UniqueId:       uniqueId,
			Currency:       payment.Currency,
			Amount:         payment.PaymentFee,
			ChannelId:      payment.ChannelId,
			PaymentId:      payment.PaymentId,
			Status:         status,
			TimelineType:   timeLineType,
		}

		result, err := dao.PaymentTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdatePaymentTimeline record insert failure %s`, err.Error())
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(id)
	} else {
		//更新
		update, err := dao.PaymentTimeline.Ctx(ctx).Data(g.Map{
			dao.PaymentTimeline.Columns().MerchantId:     payment.MerchantId,
			dao.PaymentTimeline.Columns().UserId:         payment.UserId,
			dao.PaymentTimeline.Columns().SubscriptionId: payment.SubscriptionId,
			dao.PaymentTimeline.Columns().InvoiceId:      payment.InvoiceId,
			dao.PaymentTimeline.Columns().Currency:       payment.Currency,
			dao.PaymentTimeline.Columns().Amount:         payment.PaymentFee,
			dao.PaymentTimeline.Columns().ChannelId:      payment.ChannelId,
			dao.PaymentTimeline.Columns().PaymentId:      payment.PaymentId,
			dao.PaymentTimeline.Columns().GmtModify:      gtime.Now(),
			dao.PaymentTimeline.Columns().Status:         status,
			dao.PaymentTimeline.Columns().TimelineType:   timeLineType,
		}).Where(dao.PaymentTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		rowAffected, err := update.RowsAffected()
		if rowAffected != 1 {
			return gerror.Newf("CreateOrUpdatePaymentTimeline err:%s", update)
		}
	}
	return nil
}
