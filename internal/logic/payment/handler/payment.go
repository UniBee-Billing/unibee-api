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
	"go-oversea-pay/internal/logic/channel/ro"
	handler2 "go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/internal/logic/payment/event"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

type HandlePayReq struct {
	PaymentId                        string
	ChannelPaymentIntentId           string
	ChannelPaymentId                 string
	TotalAmount                      int64
	PayStatusEnum                    consts.PayStatusEnum
	PaidTime                         *gtime.Time
	PaymentAmount                    int64
	CaptureAmount                    int64
	Reason                           string
	ChannelDefaultPaymentMethod      string
	ChannelDetailInvoiceInternalResp *ro.ChannelDetailInvoiceInternalResp
}

func HandlePayExpired(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}

	event.SaveTimeLine(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     payment.PaymentId,
		Fee:       payment.TotalAmount,
		EventType: event.Expired.Type,
		Event:     event.Expired.Desc,
		OpenApiId: payment.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Expired"),
	})

	err = handler2.UpdateInvoiceFromPayment(ctx, payment, nil)
	if err != nil {
		fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
	}

	err = CreateOrUpdatePaymentTimeline(ctx, payment, payment.PaymentId)
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
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}
	err = handler2.UpdateInvoiceFromPayment(ctx, payment, nil)
	if err != nil {
		fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
	}
	//交易事件记录
	event.SaveTimeLine(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     payment.PaymentId,
		Fee:       req.CaptureAmount,
		EventType: event.CaptureFailed.Type,
		Event:     event.CaptureFailed.Desc,
		OpenApiId: payment.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "CaptureFailed", req.ChannelPaymentIntentId),
		Message:   req.Reason,
	})
	return nil
}

func HandlePayAuthorized(ctx context.Context, payment *entity.Payment) (err error) {
	g.Log().Infof(ctx, "HandlePayAuthorized, payment=%s", utility.MarshalToJsonString(payment))
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil")
		return errors.New("支付不存在")
	}
	if payment.AuthorizeStatus == consts.AUTHORIZED {
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayAuthorized, payment.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.AUTHORIZED, dao.Payment.Columns().ChannelPaymentId: payment.ChannelPaymentId},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID, dao.Payment.Columns().AuthorizeStatus: consts.WAITING_AUTHORIZED})
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
		err = handler2.UpdateInvoiceFromPayment(ctx, payment, nil)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		event.SaveTimeLine(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.ChannelPaymentId, "Authorised"),
		})
	}

	return err

}

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePayFailure, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment null, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	if payment.Status == consts.PAY_FAILED {
		g.Log().Infof(ctx, "already failure")
		return nil
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if payment.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	var refundFee int64 = 0
	if payment.AuthorizeStatus != consts.WAITING_AUTHORIZED {
		refundFee = payment.TotalAmount
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancelld, payment.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PAY_FAILED, dao.Payment.Columns().RefundAmount: refundFee},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
		err = handler2.UpdateInvoiceFromPayment(ctx, payment, req.ChannelDetailInvoiceInternalResp)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		//交易事件记录
		event.SaveTimeLine(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       0,
			EventType: event.Cancelled.Type,
			Event:     event.Cancelled.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Cancelled"),
			Message:   req.Reason,
		})
		err = CreateOrUpdatePaymentTimeline(ctx, payment, payment.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
		}
	}
	return err
}

func HandlePaySuccess(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePaySuccess, req=%s", utility.MarshalToJsonString(req))

	if req.PaidTime == nil {
		return errors.New("invalid param PaidTime is nil")
	}
	if len(req.PaymentId) == 0 {
		return errors.New("invalid param PaymentId is nil")
	}
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)

	if payment == nil {
		g.Log().Infof(ctx, "payment not found, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if payment.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "merchantOrderNo:%s payment already success", req.PaymentId)
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaySuccess, payment.Id), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().Status:                 consts.PAY_SUCCESS,
				dao.Payment.Columns().PaidTime:               req.PaidTime,
				dao.Payment.Columns().ChannelPaymentIntentId: req.ChannelPaymentIntentId,
				dao.Payment.Columns().ChannelPaymentId:       req.ChannelPaymentId,
				dao.Payment.Columns().PaymentAmount:          req.PaymentAmount,
				dao.Payment.Columns().RefundAmount:           payment.RefundAmount},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
		//default payment method update
		if len(req.ChannelDefaultPaymentMethod) > 0 {
			_, err = dao.SubscriptionUserChannel.Ctx(ctx).Data(g.Map{
				dao.SubscriptionUserChannel.Columns().ChannelDefaultPaymentMethod: req.ChannelDefaultPaymentMethod,
			}).Where(dao.SubscriptionUserChannel.Columns().UserId, payment.UserId).Where(dao.SubscriptionUserChannel.Columns().ChannelId, payment.ChannelId).OmitNil().Update()
			if err != nil {
				g.Log().Printf(ctx, `UpdateChannelUser ChannelDefaultPaymentMethod failure %s`, err.Error())
			}
		}

		err = handler2.UpdateInvoiceFromPayment(ctx, payment, req.ChannelDetailInvoiceInternalResp)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		//try {
		//交易事件记录
		event.SaveTimeLine(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       req.PaymentAmount,
			EventType: event.Settled.Type,
			Event:     event.Settled.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Settled"),
			Message:   req.Reason,
		})
		//} catch (Exception e) {
		//	e.printStackTrace();
		//	log.info("save_event exception:{}",e.toString());
		//}
		////            mqUtil.addToMq(MqTopicEnum.PaySuccess,payment.getId());
		//// 为提高支付结果反馈速度，同步处理业务订单支付成功
		//OverseaPay queryPay = overseaPayMapper.selectById(payment.getId());
		//if (queryPay == null) {
		//	payment.setPayStatus(PayStatusEnum.PAY_SUCCESS.getCode());
		//	payment.setPaidTime(Instant.ofEpochMilli(req.getPaidTime().getTime()).atZone(ZoneId.systemDefault()).toLocalDateTime());
		//	payment.setChannelPayId(req.getChannelPayId());
		//	payment.setChannelTradeNo(req.getChannelTradeNo());
		//	payment.setReceiptFee(req.getReceiptFee());
		//	queryPay = payment;
		//}
		//com.hk.utils.BusinessWrapper result = bizOrderPayCallbackProviderFactory.getBizOrderPayCallbackServiceProvider(queryPay.getBizType()).paySuccessCallback(queryPay,req.getPaidTime());
		err := CreateOrUpdatePaymentTimeline(ctx, payment, payment.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
		}
	}
	return err
}

func HandlePaymentWebhookEvent(ctx context.Context, channelPayRo *ro.ChannelPaymentRo) error {
	one := query.GetPaymentByChannelPaymentId(ctx, channelPayRo.ChannelPaymentId)
	if channelPayRo.ChannelSubscriptionDetail != nil && channelPayRo.ChannelInvoiceDetail != nil {
		// payment for subscription
		payment, err := CreateOrUpdateSubscriptionPaymentFromChannel(ctx, channelPayRo)
		// payment not first generate from system
		if err != nil {
			return err
		}

		if len(channelPayRo.ChannelSubscriptionId) > 0 && channelPayRo.ChannelSubscriptionDetail == nil {
			return gerror.Newf("payment hook may too fast, ChannelSubscriptionDetail is nil for ChannelSubscriptionId:%s", channelPayRo.ChannelSubscriptionId)
		}
		//Subscription Payment Success
		err = handler.HandleSubscriptionPaymentUpdate(ctx, &handler.SubscriptionPaymentSuccessWebHookReq{
			Payment:                     payment,
			ChannelInvoiceDetail:        channelPayRo.ChannelInvoiceDetail,
			ChannelSubscriptionDetail:   channelPayRo.ChannelSubscriptionDetail,
			ChannelPaymentId:            channelPayRo.ChannelPaymentId,
			ChannelInvoiceId:            channelPayRo.ChannelInvoiceDetail.ChannelInvoiceId,
			ChannelSubscriptionId:       channelPayRo.ChannelInvoiceDetail.ChannelSubscriptionId,
			ChannelSubscriptionUpdateId: channelPayRo.ChannelSubscriptionUpdateId,
			Status:                      channelPayRo.ChannelSubscriptionDetail.Status,
			ChannelStatus:               channelPayRo.ChannelInvoiceDetail.ChannelStatus,
			Data:                        channelPayRo.ChannelSubscriptionDetail.Data,
			ChannelItemData:             channelPayRo.ChannelSubscriptionDetail.ChannelItemData,
			CancelAtPeriodEnd:           channelPayRo.ChannelSubscriptionDetail.CancelAtPeriodEnd,
			CurrentPeriodEnd:            channelPayRo.ChannelSubscriptionDetail.CurrentPeriodEnd,
			CurrentPeriodStart:          channelPayRo.ChannelSubscriptionDetail.CurrentPeriodStart,
			TrialEnd:                    channelPayRo.ChannelSubscriptionDetail.TrialEnd,
		})
		err = handler2.UpdateInvoiceFromPayment(ctx, payment, channelPayRo.ChannelInvoiceDetail)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
	} else if one != nil {
		// one-time payment
		if channelPayRo.Status == consts.PAY_SUCCESS {
			err := HandlePaySuccess(ctx, &HandlePayReq{
				PaymentId:                        one.PaymentId,
				ChannelPaymentIntentId:           channelPayRo.ChannelPaymentId,
				ChannelPaymentId:                 channelPayRo.ChannelPaymentId,
				TotalAmount:                      channelPayRo.TotalAmount,
				PayStatusEnum:                    consts.PAY_SUCCESS,
				PaidTime:                         channelPayRo.PayTime,
				PaymentAmount:                    channelPayRo.PaymentAmount,
				CaptureAmount:                    0,
				Reason:                           channelPayRo.Reason,
				ChannelDefaultPaymentMethod:      "",
				ChannelDetailInvoiceInternalResp: channelPayRo.ChannelInvoiceDetail,
			})
			if err != nil {
				return err
			}
			if consts.PendingSubUpdateEffectImmediateWithOutChannel {
				// better use redis mq to trace payment
				if len(channelPayRo.ChannelSubscriptionUpdateId) > 0 {
					eiPendingSubUpdate := query.GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx, channelPayRo.ChannelSubscriptionUpdateId)
					if eiPendingSubUpdate != nil {
						//更新单支付成功, EffectImmediate=true 需要用户 3DS 验证等场景
						sub := query.GetSubscriptionBySubscriptionId(ctx, eiPendingSubUpdate.SubscriptionId)
						_, err := handler.FinishPendingUpdateForSubscription(ctx, sub, eiPendingSubUpdate)
						if err != nil {
							return err
						}
					}
				}
			}
		} else if channelPayRo.Status == consts.PAY_FAILED {
			err := HandlePayFailure(ctx, &HandlePayReq{
				PaymentId:                        one.PaymentId,
				ChannelPaymentIntentId:           channelPayRo.ChannelPaymentId,
				ChannelPaymentId:                 channelPayRo.ChannelPaymentId,
				TotalAmount:                      channelPayRo.TotalAmount,
				PayStatusEnum:                    consts.PAY_FAILED,
				PaidTime:                         channelPayRo.PayTime,
				PaymentAmount:                    channelPayRo.PaymentAmount,
				CaptureAmount:                    0,
				Reason:                           channelPayRo.Reason,
				ChannelDefaultPaymentMethod:      "",
				ChannelDetailInvoiceInternalResp: channelPayRo.ChannelInvoiceDetail,
			})
			if err != nil {
				return err
			}
		}
	} else {
		return gerror.Newf("Invalid Payment Type ChannelPaymentId:%s ChannelInvoiceId:%s ChannelSubscriptionId:%s", channelPayRo.ChannelPaymentId, channelPayRo.ChannelInvoiceId, channelPayRo.ChannelSubscriptionId)
	}

	return nil
}

func CreateOrUpdateSubscriptionPaymentFromChannel(ctx context.Context, channelPayRo *ro.ChannelPaymentRo) (*entity.Payment, error) {
	utility.Assert(len(channelPayRo.UniqueId) > 0, "uniqueId invalid")
	channelUser := query.GetUserChannelByChannelUserId(ctx, channelPayRo.ChannelUserId, channelPayRo.ChannelId)
	utility.Assert(channelUser != nil, "channelUser not found")
	var subscriptionId string
	var invoiceId string
	var countryCode string
	if channelPayRo.ChannelSubscriptionDetail != nil {
		//From Sub Create Pay or Sub Update Pay
		sub := query.GetSubscriptionByChannelSubscriptionId(ctx, channelPayRo.ChannelSubscriptionDetail.ChannelSubscriptionId)
		if sub != nil {
			subscriptionId = sub.SubscriptionId
			countryCode = sub.CountryCode
		}
	}
	if channelPayRo.ChannelInvoiceDetail != nil {
		//From Invoice Pay
		invoice := query.GetInvoiceByChannelInvoiceId(ctx, channelPayRo.ChannelInvoiceId)
		if invoice != nil {
			invoiceId = invoice.InvoiceId
		}
	}
	one := query.GetPaymentByChannelUniqueId(ctx, channelPayRo.UniqueId)
	if one == nil {
		//创建
		one = &entity.Payment{
			BizType:                consts.BIZ_TYPE_SUBSCRIPTION,
			MerchantId:             channelPayRo.MerchantId,
			UserId:                 channelUser.UserId,
			CountryCode:            countryCode,
			PaymentId:              utility.CreatePaymentId(),
			Currency:               channelPayRo.Currency,
			TotalAmount:            channelPayRo.TotalAmount,
			PaymentAmount:          channelPayRo.PaymentAmount,
			BalanceAmount:          channelPayRo.BalanceAmount,
			BalanceStart:           channelPayRo.BalanceStart,
			BalanceEnd:             channelPayRo.BalanceEnd,
			Status:                 channelPayRo.Status,
			AuthorizeStatus:        channelPayRo.CaptureStatus,
			ChannelId:              channelPayRo.ChannelId,
			ChannelPaymentIntentId: channelPayRo.ChannelPaymentId,
			ChannelPaymentId:       channelPayRo.ChannelPaymentId,
			CreateTime:             channelPayRo.CreateTime,
			CancelTime:             channelPayRo.CancelTime,
			PaidTime:               channelPayRo.PayTime,
			PaymentData:            channelPayRo.CancelReason,
			UniqueId:               channelPayRo.UniqueId,
			SubscriptionId:         subscriptionId,
			InvoiceId:              invoiceId,
		}
		result, err := dao.Payment.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdateSubscriptionPaymentFromChannel record insert failure %s`, err.Error())
			return nil, err
		}
		id, _ := result.LastInsertId()
		one.Id = id
	} else {
		//更新
		_, err := dao.Payment.Ctx(ctx).Data(g.Map{
			dao.Payment.Columns().BizType:                consts.BIZ_TYPE_SUBSCRIPTION,
			dao.Payment.Columns().MerchantId:             channelPayRo.MerchantId,
			dao.Payment.Columns().UserId:                 channelUser.UserId,
			dao.Payment.Columns().CountryCode:            countryCode,
			dao.Payment.Columns().Currency:               channelPayRo.Currency,
			dao.Payment.Columns().TotalAmount:            channelPayRo.TotalAmount,
			dao.Payment.Columns().PaymentAmount:          channelPayRo.PaymentAmount,
			dao.Payment.Columns().BalanceAmount:          channelPayRo.BalanceAmount,
			dao.Payment.Columns().BalanceStart:           channelPayRo.BalanceStart,
			dao.Payment.Columns().BalanceEnd:             channelPayRo.BalanceEnd,
			dao.Payment.Columns().Status:                 channelPayRo.Status,
			dao.Payment.Columns().AuthorizeStatus:        channelPayRo.CaptureStatus,
			dao.Payment.Columns().ChannelId:              channelPayRo.ChannelId,
			dao.Payment.Columns().ChannelPaymentIntentId: channelPayRo.ChannelPaymentId,
			dao.Payment.Columns().CreateTime:             channelPayRo.CreateTime,
			dao.Payment.Columns().CancelTime:             channelPayRo.CancelTime,
			dao.Payment.Columns().PaidTime:               channelPayRo.PayTime,
			dao.Payment.Columns().PaymentData:            channelPayRo.CancelReason,
			dao.Payment.Columns().SubscriptionId:         subscriptionId,
			dao.Payment.Columns().InvoiceId:              invoiceId,
			dao.Invoice.Columns().GmtModify:              gtime.Now(),
		}).Where(dao.Payment.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return nil, err
		}
	}
	err := CreateOrUpdatePaymentTimeline(ctx, one, one.PaymentId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
	}
	return query.GetPaymentByChannelUniqueId(ctx, channelPayRo.UniqueId), nil
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
	if payment.TotalAmount > 0 {
		timeLineType = 0
	} else if payment.TotalAmount < 0 {
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
			TotalAmount:    payment.TotalAmount,
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
		_, err := dao.PaymentTimeline.Ctx(ctx).Data(g.Map{
			dao.PaymentTimeline.Columns().MerchantId:     payment.MerchantId,
			dao.PaymentTimeline.Columns().UserId:         payment.UserId,
			dao.PaymentTimeline.Columns().SubscriptionId: payment.SubscriptionId,
			dao.PaymentTimeline.Columns().InvoiceId:      payment.InvoiceId,
			dao.PaymentTimeline.Columns().Currency:       payment.Currency,
			dao.PaymentTimeline.Columns().TotalAmount:    payment.TotalAmount,
			dao.PaymentTimeline.Columns().ChannelId:      payment.ChannelId,
			dao.PaymentTimeline.Columns().PaymentId:      payment.PaymentId,
			dao.PaymentTimeline.Columns().GmtModify:      gtime.Now(),
			dao.PaymentTimeline.Columns().Status:         status,
			dao.PaymentTimeline.Columns().TimelineType:   timeLineType,
		}).Where(dao.PaymentTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			return err
		}
		//rowAffected, err := update.RowsAffected()
		//if rowAffected != 1 {
		//	return gerror.Newf("CreateOrUpdatePaymentTimeline err:%s", update)
		//}
	}
	return nil
}
