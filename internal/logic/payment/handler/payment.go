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
	handler2 "go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/internal/logic/payment/callback"
	"go-oversea-pay/internal/logic/payment/event"
	"go-oversea-pay/internal/logic/subscription/handler"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

type HandlePayReq struct {
	PaymentId                        string
	GatewayPaymentIntentId           string
	GatewayPaymentId                 string
	TotalAmount                      int64
	PayStatusEnum                    consts.PayStatusEnum
	PaidTime                         *gtime.Time
	PaymentAmount                    int64
	CaptureAmount                    int64
	Reason                           string
	ChannelDefaultPaymentMethod      string
	ChannelDetailInvoiceInternalResp *ro.GatewayDetailInvoiceInternalResp
}

func HandlePayExpired(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayExpired, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil, paymentId=%s", req.PaymentId)
		return errors.New("支付不存在")
	}

	event.SaveEvent(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     payment.PaymentId,
		Fee:       payment.TotalAmount,
		EventType: event.Expired.Type,
		Event:     event.Expired.Desc,
		OpenApiId: payment.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Expired"),
	})

	_, err = handler2.UpdateInvoiceFromPayment(ctx, payment)
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
	_, err = handler2.UpdateInvoiceFromPayment(ctx, payment)
	if err != nil {
		fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
	}
	//交易事件记录
	event.SaveEvent(ctx, entity.PaymentEvent{
		BizType:   0,
		BizId:     payment.PaymentId,
		Fee:       req.CaptureAmount,
		EventType: event.CaptureFailed.Type,
		Event:     event.CaptureFailed.Desc,
		OpenApiId: payment.OpenApiId,
		UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "CaptureFailed", req.GatewayPaymentIntentId),
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

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayAuthorized, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.AUTHORIZED, dao.Payment.Columns().GatewayPaymentId: payment.GatewayPaymentId},
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
		_, err = handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.GatewayPaymentId, "Authorised"),
		})
	}
	return err
}

func HandlePayNeedAuthorized(ctx context.Context, payment *entity.Payment, authorizeReason string, paymentData string) (err error) {
	g.Log().Infof(ctx, "HandlePayNeedAuthorized, payment=%s", utility.MarshalToJsonString(payment))
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil")
		return errors.New("payment not found")
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayAuthorized, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().AuthorizeStatus:  consts.WAITING_AUTHORIZED,
				dao.Payment.Columns().AuthorizeReason:  authorizeReason,
				dao.Payment.Columns().PaymentData:      paymentData,
				dao.Payment.Columns().GatewayPaymentId: payment.GatewayPaymentId},
				g.Map{dao.Payment.Columns().Id: payment.Id,
					dao.Payment.Columns().Status: consts.TO_BE_PAID})
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
	g.Log().Infof(ctx, "HandlePayNeedAuthorized sendResult err=%s", err)
	if err == nil {
		payment = query.GetPaymentByPaymentId(ctx, payment.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentNeedAuthorisedCallback(ctx, payment, invoice)

		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.GatewayPaymentId, "NeedAuthorised"),
		})
	}
	return err
}

func HandlePayCancel(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "HandlePayCancel, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment null, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	if payment.Status == consts.PAY_CANCEL || payment.Status == consts.PAY_FAILED {
		g.Log().Infof(ctx, "already cancel or failure")
		return nil
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if payment.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancel, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PAY_CANCEL, dao.Payment.Columns().CancelTime: gtime.Now(), dao.Payment.Columns().FailureReason: req.Reason},
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
			payment.Status = consts.PAY_CANCEL
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

	g.Log().Infof(ctx, "HandlePayCancel sendResult err=%s", err)
	if err == nil {
		payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}

		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentCancelCallback(ctx, payment, invoice)
		//交易事件记录
		event.SaveEvent(ctx, entity.PaymentEvent{
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

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "handlePayFailure, req=%s", utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment null, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	if payment.Status == consts.PAY_CANCEL || payment.Status == consts.PAY_FAILED {
		g.Log().Infof(ctx, "already cancel or failure")
		return nil
	}

	// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	if payment.Status == consts.PAY_SUCCESS {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPayCancel, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PAY_FAILED, dao.Payment.Columns().CancelTime: gtime.Now(), dao.Payment.Columns().FailureReason: req.Reason},
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
			payment.Status = consts.PAY_FAILED
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
		payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}

		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentFailureCallback(ctx, payment, invoice)
		//交易事件记录
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       0,
			EventType: event.Cancelled.Type,
			Event:     event.Cancelled.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Failed"),
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

	//// 支付宝存在 TRADE_FINISHED 交易完结  https://opendocs.alipay.com/open/02ekfj?ref=api
	//if payment.Status == consts.PAY_SUCCESS {
	//	g.Log().Infof(ctx, "merchantOrderNo:%s payment already success", req.PaymentId)
	//	return nil
	//}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaySuccess, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().Status:                 consts.PAY_SUCCESS,
				dao.Payment.Columns().PaidTime:               req.PaidTime,
				dao.Payment.Columns().GatewayPaymentIntentId: req.GatewayPaymentIntentId,
				dao.Payment.Columns().GatewayPaymentId:       req.GatewayPaymentId,
				dao.Payment.Columns().GatewayPaymentMethod:   req.ChannelDefaultPaymentMethod,
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
			payment.Status = consts.PAY_SUCCESS
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
		payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}

		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentSuccessCallback(ctx, payment, invoice)

		//default payment method update
		if len(req.ChannelDefaultPaymentMethod) > 0 {
			_ = SaveChannelUserDefaultPaymentMethod(ctx, req, err, payment)
		}

		//try {
		//交易事件记录
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       req.PaymentAmount,
			EventType: event.Settled.Type,
			Event:     event.Settled.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s", payment.PaymentId, "Settled"),
			Message:   req.Reason,
		})

		err = CreateOrUpdatePaymentTimeline(ctx, payment, payment.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimeline error %s`, err.Error())
		}
	}
	return err
}

func SaveChannelUserDefaultPaymentMethod(ctx context.Context, req *HandlePayReq, err error, payment *entity.Payment) error {
	_, err = dao.GatewayUser.Ctx(ctx).Data(g.Map{
		dao.GatewayUser.Columns().GatewayDefaultPaymentMethod: req.ChannelDefaultPaymentMethod,
	}).Where(dao.GatewayUser.Columns().UserId, payment.UserId).Where(dao.GatewayUser.Columns().GatewayId, payment.GatewayId).OmitNil().Update()
	if err != nil {
		g.Log().Printf(ctx, `SaveChannelUserDefaultPaymentMethod GatewayDefaultPaymentMethod failure %s`, err.Error())
	}
	return err
}

func HandlePaymentWebhookEvent(ctx context.Context, gatewayPaymentRo *ro.GatewayPaymentRo) error {
	one := query.GetPaymentByGatewayPaymentId(ctx, gatewayPaymentRo.GatewayPaymentId)
	if gatewayPaymentRo.GatewaySubscriptionDetail != nil && gatewayPaymentRo.GatewayInvoiceDetail != nil {
		// payment for subscription
		payment, err := CreateOrUpdateSubscriptionPaymentFromChannel(ctx, gatewayPaymentRo)
		// payment not first generate from system
		if err != nil {
			return err
		}

		if len(gatewayPaymentRo.GatewaySubscriptionId) > 0 && gatewayPaymentRo.GatewaySubscriptionDetail == nil {
			return gerror.Newf("payment hook may too fast, GatewaySubscriptionDetail is nil for GatewaySubscriptionId:%s", gatewayPaymentRo.GatewaySubscriptionId)
		}
		//Subscription Payment Success
		err = handler.HandleSubscriptionPaymentUpdate(ctx, &handler.SubscriptionPaymentSuccessWebHookReq{
			Payment:                     payment,
			GatewayInvoiceDetail:        gatewayPaymentRo.GatewayInvoiceDetail,
			GatewaySubscriptionDetail:   gatewayPaymentRo.GatewaySubscriptionDetail,
			GatewayPaymentId:            gatewayPaymentRo.GatewayPaymentId,
			GatewayInvoiceId:            gatewayPaymentRo.GatewayInvoiceDetail.GatewayInvoiceId,
			GatewaySubscriptionId:       gatewayPaymentRo.GatewayInvoiceDetail.GatewaySubscriptionId,
			GatewaySubscriptionUpdateId: gatewayPaymentRo.GatewaySubscriptionUpdateId,
			Status:                      gatewayPaymentRo.GatewaySubscriptionDetail.Status,
			GatewayStatus:               gatewayPaymentRo.GatewayInvoiceDetail.GatewayStatus,
			Data:                        gatewayPaymentRo.GatewaySubscriptionDetail.Data,
			GatewayItemData:             gatewayPaymentRo.GatewaySubscriptionDetail.GatewayItemData,
			CancelAtPeriodEnd:           gatewayPaymentRo.GatewaySubscriptionDetail.CancelAtPeriodEnd,
			CurrentPeriodEnd:            gatewayPaymentRo.GatewaySubscriptionDetail.CurrentPeriodEnd,
			CurrentPeriodStart:          gatewayPaymentRo.GatewaySubscriptionDetail.CurrentPeriodStart,
			TrialEnd:                    gatewayPaymentRo.GatewaySubscriptionDetail.TrialEnd,
		})
		_, err = handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s`, err.Error())
		}
	} else if one != nil {
		// one-time payment
		if gatewayPaymentRo.Status == consts.PAY_SUCCESS {
			err := HandlePaySuccess(ctx, &HandlePayReq{
				PaymentId:                        one.PaymentId,
				GatewayPaymentIntentId:           gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:                 gatewayPaymentRo.GatewayPaymentId,
				TotalAmount:                      gatewayPaymentRo.TotalAmount,
				PayStatusEnum:                    consts.PAY_SUCCESS,
				PaidTime:                         gatewayPaymentRo.PayTime,
				PaymentAmount:                    gatewayPaymentRo.PaymentAmount,
				CaptureAmount:                    0,
				Reason:                           gatewayPaymentRo.Reason,
				ChannelDefaultPaymentMethod:      gatewayPaymentRo.GatewayPaymentMethod,
				ChannelDetailInvoiceInternalResp: gatewayPaymentRo.GatewayInvoiceDetail,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.Status == consts.PAY_FAILED {
			err := HandlePayFailure(ctx, &HandlePayReq{
				PaymentId:                        one.PaymentId,
				GatewayPaymentIntentId:           gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:                 gatewayPaymentRo.GatewayPaymentId,
				TotalAmount:                      gatewayPaymentRo.TotalAmount,
				PayStatusEnum:                    consts.PAY_FAILED,
				PaidTime:                         gatewayPaymentRo.PayTime,
				PaymentAmount:                    gatewayPaymentRo.PaymentAmount,
				CaptureAmount:                    0,
				Reason:                           gatewayPaymentRo.Reason,
				ChannelDetailInvoiceInternalResp: gatewayPaymentRo.GatewayInvoiceDetail,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.Status == consts.PAY_CANCEL {
			err := HandlePayCancel(ctx, &HandlePayReq{
				PaymentId:                        one.PaymentId,
				GatewayPaymentIntentId:           gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:                 gatewayPaymentRo.GatewayPaymentId,
				TotalAmount:                      gatewayPaymentRo.TotalAmount,
				PayStatusEnum:                    consts.PAY_CANCEL,
				PaidTime:                         gatewayPaymentRo.PayTime,
				PaymentAmount:                    gatewayPaymentRo.PaymentAmount,
				CaptureAmount:                    0,
				Reason:                           gatewayPaymentRo.Reason,
				ChannelDetailInvoiceInternalResp: gatewayPaymentRo.GatewayInvoiceDetail,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.AuthorizeStatus == consts.WAITING_AUTHORIZED {
			err := HandlePayNeedAuthorized(ctx, one, gatewayPaymentRo.AuthorizeReason, gatewayPaymentRo.PaymentData)
			if err != nil {
				return err
			}
		}
	} else {
		return gerror.Newf("Payment Not Match Or Not Found GatewayPaymentId:%s GatewayInvoiceId:%s GatewaySubscriptionId:%s", gatewayPaymentRo.GatewayPaymentId, gatewayPaymentRo.GatewayInvoiceId, gatewayPaymentRo.GatewaySubscriptionId)
	}

	return nil
}

func CreateOrUpdateSubscriptionPaymentFromChannel(ctx context.Context, gatewayPaymentRo *ro.GatewayPaymentRo) (*entity.Payment, error) {
	utility.Assert(len(gatewayPaymentRo.UniqueId) > 0, "uniqueId invalid")
	gatewayUser := query.GetGatewayUserByGatewayUserId(ctx, gatewayPaymentRo.GatewayUserId, gatewayPaymentRo.GatewayId)
	utility.Assert(gatewayUser != nil, "gatewayUser not found")
	var subscriptionId string
	var invoiceId string
	var countryCode string
	if gatewayPaymentRo.GatewaySubscriptionDetail != nil {
		//From Sub Create Pay or Sub Update Pay
		sub := query.GetSubscriptionByGatewaySubscriptionId(ctx, gatewayPaymentRo.GatewaySubscriptionDetail.GatewaySubscriptionId)
		if sub != nil {
			subscriptionId = sub.SubscriptionId
			countryCode = sub.CountryCode
		}
	}
	if gatewayPaymentRo.GatewayInvoiceDetail != nil {
		//From Invoice Pay
		invoice := query.GetInvoiceByGatewayInvoiceId(ctx, gatewayPaymentRo.GatewayInvoiceId)
		if invoice != nil {
			invoiceId = invoice.InvoiceId
		}
	}
	one := query.GetPaymentByGatewayUniqueId(ctx, gatewayPaymentRo.UniqueId)
	if one == nil {
		//创建
		one = &entity.Payment{
			BizType:                consts.BIZ_TYPE_SUBSCRIPTION,
			MerchantId:             gatewayPaymentRo.MerchantId,
			UserId:                 gatewayUser.UserId,
			CountryCode:            countryCode,
			PaymentId:              utility.CreatePaymentId(),
			Currency:               gatewayPaymentRo.Currency,
			TotalAmount:            gatewayPaymentRo.TotalAmount,
			PaymentAmount:          gatewayPaymentRo.PaymentAmount,
			BalanceAmount:          gatewayPaymentRo.BalanceAmount,
			BalanceStart:           gatewayPaymentRo.BalanceStart,
			BalanceEnd:             gatewayPaymentRo.BalanceEnd,
			Status:                 gatewayPaymentRo.Status,
			AuthorizeStatus:        gatewayPaymentRo.AuthorizeStatus,
			GatewayId:              gatewayPaymentRo.GatewayId,
			GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
			GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
			CreateTime:             gatewayPaymentRo.CreateTime,
			CancelTime:             gatewayPaymentRo.CancelTime,
			PaidTime:               gatewayPaymentRo.PayTime,
			FailureReason:          gatewayPaymentRo.CancelReason,
			PaymentData:            gatewayPaymentRo.PaymentData,
			AuthorizeReason:        gatewayPaymentRo.AuthorizeReason,
			UniqueId:               gatewayPaymentRo.UniqueId,
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
			dao.Payment.Columns().MerchantId:             gatewayPaymentRo.MerchantId,
			dao.Payment.Columns().UserId:                 gatewayUser.UserId,
			dao.Payment.Columns().CountryCode:            countryCode,
			dao.Payment.Columns().Currency:               gatewayPaymentRo.Currency,
			dao.Payment.Columns().TotalAmount:            gatewayPaymentRo.TotalAmount,
			dao.Payment.Columns().PaymentAmount:          gatewayPaymentRo.PaymentAmount,
			dao.Payment.Columns().BalanceAmount:          gatewayPaymentRo.BalanceAmount,
			dao.Payment.Columns().BalanceStart:           gatewayPaymentRo.BalanceStart,
			dao.Payment.Columns().BalanceEnd:             gatewayPaymentRo.BalanceEnd,
			dao.Payment.Columns().Status:                 gatewayPaymentRo.Status,
			dao.Payment.Columns().AuthorizeStatus:        gatewayPaymentRo.AuthorizeStatus,
			dao.Payment.Columns().GatewayId:              gatewayPaymentRo.GatewayId,
			dao.Payment.Columns().GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
			dao.Payment.Columns().CreateTime:             gatewayPaymentRo.CreateTime,
			dao.Payment.Columns().CancelTime:             gatewayPaymentRo.CancelTime,
			dao.Payment.Columns().PaidTime:               gatewayPaymentRo.PayTime,
			dao.Payment.Columns().FailureReason:          gatewayPaymentRo.CancelReason,
			dao.Payment.Columns().PaymentData:            gatewayPaymentRo.PaymentData,
			dao.Payment.Columns().AuthorizeReason:        gatewayPaymentRo.AuthorizeReason,
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
	return query.GetPaymentByGatewayUniqueId(ctx, gatewayPaymentRo.UniqueId), nil
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
			GatewayId:      payment.GatewayId,
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
			dao.PaymentTimeline.Columns().GatewayId:      payment.GatewayId,
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
