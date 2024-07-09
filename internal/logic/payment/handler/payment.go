package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	event2 "unibee/internal/consumer/webhook/event"
	payment2 "unibee/internal/consumer/webhook/payment"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/gateway/api"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/logic/payment/event"
	"unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type HandlePayReq struct {
	PaymentId              string
	GatewayPaymentIntentId string
	GatewayPaymentId       string
	GatewayUserId          string
	TotalAmount            int64
	PayStatusEnum          consts.PaymentStatusEnum
	PaidTime               *gtime.Time
	PaymentAmount          int64
	CaptureAmount          int64
	Reason                 string
	GatewayPaymentMethod   string
	PaymentCode            string
}

func UpdatePaymentLastGatewayError(ctx context.Context, paymentId string, lastErr string) {
	_, _ = dao.Payment.Ctx(ctx).Data(g.Map{
		dao.Payment.Columns().LastError: lastErr,
	}).Where(dao.Payment.Columns().PaymentId, paymentId).OmitNil().Update()
}

func HandlePayExpired(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "%s-HandlePayExpired, req=%s", req.PaymentId, utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
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
		fmt.Printf(`UpdateInvoiceFromPayment error %s\n`, err.Error())
	}

	err = CreateOrUpdatePaymentTimelineForPayment(ctx, payment, payment.PaymentId)
	if err != nil {
		fmt.Printf(`CreateOrUpdatePaymentTimelineForPayment error %s`, err.Error())
	}

	return HandlePayFailure(ctx, &HandlePayReq{
		PaymentId:     req.PaymentId,
		PayStatusEnum: consts.PaymentFailed,
		Reason:        "FailedByExpire",
	})
}

func HandleCaptureFailed(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "%s-HandlePayExpired, req=%s", req.PaymentId, utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment is nil, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	_, err = handler2.UpdateInvoiceFromPayment(ctx, payment)
	if err != nil {
		g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
	}

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
	if payment == nil {
		g.Log().Infof(ctx, "HandlePayAuthorized payment is nil")
		return errors.New("payment not found")
	}
	g.Log().Infof(ctx, "%s-HandlePayAuthorized, payment=%s", payment.PaymentId, utility.MarshalToJsonString(payment))
	if payment.AuthorizeStatus == consts.Authorized {
		return nil
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaymentAuthorized, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().AuthorizeStatus: consts.Authorized, dao.Payment.Columns().GatewayPaymentId: payment.GatewayPaymentId},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.PaymentCreated, dao.Payment.Columns().AuthorizeStatus: consts.WaitingAuthorized})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
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
			g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
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

func HandlePayNeedAuthorized(ctx context.Context, payment *entity.Payment, authorizeReason string, paymentData string, paymentCode string) (err error) {
	if payment == nil {
		g.Log().Infof(ctx, "HandlePayNeedAuthorized payment is nil")
		return errors.New("payment not found")
	}
	g.Log().Infof(ctx, "%s-HandlePayNeedAuthorized, payment=%s", payment.PaymentId, utility.MarshalToJsonString(payment))

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaymentAuthorized, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().AuthorizeStatus:  consts.WaitingAuthorized,
				dao.Payment.Columns().AuthorizeReason:  authorizeReason,
				dao.Payment.Columns().PaymentData:      paymentData,
				dao.Payment.Columns().Code:             paymentCode,
				dao.Payment.Columns().GatewayPaymentId: payment.GatewayPaymentId},
				g.Map{dao.Payment.Columns().Id: payment.Id,
					dao.Payment.Columns().Status: consts.PaymentCreated})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
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
			g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
		}
		payment2.SendPaymentWebhookBackground(payment.PaymentId, event2.UNIBEE_WEBHOOK_EVENT_PAYMENT_NEEDAUTHORISED)
		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentNeedAuthorisedCallback(ctx, payment, invoice)
		if err != nil {
			fmt.Printf(`UpdateInvoiceFromPayment error %s\n`, err.Error())
		}
		event.SaveEvent(ctx, entity.PaymentEvent{
			BizType:   0,
			BizId:     payment.PaymentId,
			Fee:       payment.TotalAmount,
			EventType: event.Authorised.Type,
			Event:     event.Authorised.Desc,
			OpenApiId: payment.OpenApiId,
			UniqueNo:  fmt.Sprintf("%s_%s_%s", payment.PaymentId, "NeedAuthorised", payment.GatewayPaymentId),
		})
	}
	return err
}

func HandlePayCancel(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "%s-HandlePayCancel, req=%s", req.PaymentId, utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment null, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	if payment.Status == consts.PaymentCancelled || payment.Status == consts.PaymentFailed {
		g.Log().Infof(ctx, "already cancel or failure")
		return nil
	}

	if payment.Status == consts.PaymentSuccess {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	if len(req.Reason) == 0 {
		req.Reason = payment.FailureReason
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaymentCancel, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PaymentCancelled, dao.Payment.Columns().CancelTime: gtime.Now().Timestamp(), dao.Payment.Columns().FailureReason: req.Reason},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.PaymentCreated})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				return err
			}
			payment.Status = consts.PaymentCancelled
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
		payment = query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
		}

		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentCancelCallback(ctx, payment, invoice)

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
		err = CreateOrUpdatePaymentTimelineForPayment(ctx, payment, payment.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineForPayment error %s`, err.Error())
		}
	}
	return err
}

func HandlePayFailure(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "%s-handlePayFailure, req=%s", req.PaymentId, utility.MarshalToJsonString(req))
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	if payment == nil {
		g.Log().Infof(ctx, "payment null, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}
	if payment.Status == consts.PaymentCancelled || payment.Status == consts.PaymentFailed {
		g.Log().Infof(ctx, "already cancel or failure")
		return nil
	}

	if payment.Status == consts.PaymentSuccess {
		g.Log().Infof(ctx, "payment already success")
		return errors.New("payment already success")
	}

	if len(req.Reason) == 0 {
		req.Reason = payment.FailureReason
	}

	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaymentCancel, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			result, err := transaction.Update(dao.Payment.Table(), g.Map{dao.Payment.Columns().Status: consts.PaymentFailed, dao.Payment.Columns().CancelTime: gtime.Now().Timestamp(), dao.Payment.Columns().FailureReason: req.Reason},
				g.Map{dao.Payment.Columns().Id: payment.Id, dao.Payment.Columns().Status: consts.PaymentCreated})
			if err != nil || result == nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil || affected != 1 {
				return err
			}
			payment.Status = consts.PaymentFailed
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
		payment = query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
		}

		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentFailureCallback(ctx, payment, invoice)
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
		err = CreateOrUpdatePaymentTimelineForPayment(ctx, payment, payment.PaymentId)
		if err != nil {
			fmt.Printf(`CreateOrUpdatePaymentTimelineForPayment error %s`, err.Error())
		}
	}
	return err
}

func HandlePaySuccess(ctx context.Context, req *HandlePayReq) (err error) {
	g.Log().Infof(ctx, "%s-handlePaySuccess, req=%s", req.PaymentId, utility.MarshalToJsonString(req))
	if len(req.PaymentId) == 0 {
		return errors.New("invalid param PaymentId is nil")
	}
	payment := query.GetPaymentByPaymentId(ctx, req.PaymentId)

	if payment == nil {
		g.Log().Infof(ctx, "payment not found, paymentId=%s", req.PaymentId)
		return errors.New("payment not found")
	}

	if payment.Status == consts.PaymentSuccess {
		g.Log().Infof(ctx, "payment already success, paymentId=%s", req.PaymentId)
		if len(req.PaymentCode) == 0 {
			req.PaymentCode = payment.Code
		}
		if len(req.GatewayPaymentMethod) == 0 {
			req.GatewayPaymentMethod = payment.GatewayPaymentMethod
		}
		_, _ = dao.Payment.Ctx(ctx).Data(g.Map{
			dao.Payment.Columns().Code:                 req.PaymentCode,
			dao.Payment.Columns().GatewayPaymentMethod: req.GatewayPaymentMethod,
		}).Where(dao.Payment.Columns().PaymentId, payment.PaymentId).OmitNil().Update()
		return nil
	}

	var paidAt = gtime.Now().Timestamp()
	if req.PaidTime != nil {
		paidAt = req.PaidTime.Timestamp()
	}
	if len(req.PaymentCode) == 0 {
		req.PaymentCode = payment.Code
	}
	if len(req.GatewayPaymentMethod) == 0 {
		req.GatewayPaymentMethod = payment.GatewayPaymentMethod
	}
	_, err = redismq.SendTransaction(redismq.NewRedisMQMessage(redismqcmd.TopicPaymentSuccess, payment.PaymentId), func(messageToSend *redismq.Message) (redismq.TransactionStatus, error) {
		err = dao.Payment.DB().Transaction(ctx, func(ctx context.Context, transaction gdb.TX) error {
			_, err = transaction.Update(dao.Payment.Table(), g.Map{
				dao.Payment.Columns().Status:                 consts.PaymentSuccess,
				dao.Payment.Columns().PaidTime:               paidAt,
				dao.Payment.Columns().GatewayPaymentIntentId: req.GatewayPaymentIntentId,
				dao.Payment.Columns().GatewayPaymentId:       req.GatewayPaymentId,
				dao.Payment.Columns().GatewayPaymentMethod:   req.GatewayPaymentMethod,
				dao.Payment.Columns().PaymentAmount:          req.PaymentAmount,
				dao.Payment.Columns().Code:                   req.PaymentCode,
			},
				g.Map{dao.Payment.Columns().Id: payment.Id})
			if err != nil {
				return err
			}
			payment.Status = consts.PaymentSuccess
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
		payment = query.GetPaymentByPaymentId(ctx, req.PaymentId)
		invoice, err := handler2.UpdateInvoiceFromPayment(ctx, payment)
		if err != nil {
			g.Log().Errorf(ctx, `UpdateInvoiceFromPayment error %s\n`, err.Error())
		}
		if len(req.GatewayUserId) > 0 {
			_, _ = query.CreateOrUpdateGatewayUser(ctx, payment.UserId, payment.GatewayId, req.GatewayUserId, req.GatewayPaymentMethod)
		}
		callback.GetPaymentCallbackServiceProvider(ctx, payment.BizType).PaymentSuccessCallback(ctx, payment, invoice)
		{
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
			err = CreateOrUpdatePaymentTimelineForPayment(ctx, payment, payment.PaymentId)
			if err != nil {
				g.Log().Errorf(ctx, `CreateOrUpdatePaymentTimelineForPayment error %s`, err.Error())
			}
		}
		{
			//default payment method update
			if len(req.GatewayPaymentMethod) > 0 {
				_ = SaveChannelUserDefaultPaymentMethod(ctx, req, err, payment)
				user.UpdateUserDefaultGatewayPaymentMethod(ctx, payment.UserId, payment.GatewayId, req.GatewayPaymentMethod)
			}
			gatewayUser := query.GetGatewayUser(ctx, payment.UserId, payment.GatewayId)
			gateway := query.GetGatewayById(ctx, payment.GatewayId)
			if gatewayUser != nil && gateway != nil && len(payment.GatewayPaymentMethod) > 0 {
				_, _ = query.CreateOrUpdateGatewayUser(ctx, payment.UserId, payment.GatewayId, gatewayUser.GatewayUserId, payment.GatewayPaymentMethod)
				_, _ = api.GetGatewayServiceProvider(ctx, gatewayUser.GatewayId).GatewayUserAttachPaymentMethodQuery(ctx, gateway, gatewayUser.UserId, payment.GatewayPaymentMethod)
			}
		}
	}
	return err
}

func SaveChannelUserDefaultPaymentMethod(ctx context.Context, req *HandlePayReq, err error, payment *entity.Payment) error {
	if len(payment.GatewayPaymentMethod) == 0 {
		return nil
	}
	_, err = dao.GatewayUser.Ctx(ctx).Data(g.Map{
		dao.GatewayUser.Columns().GatewayDefaultPaymentMethod: req.GatewayPaymentMethod,
	}).Where(dao.GatewayUser.Columns().UserId, payment.UserId).Where(dao.GatewayUser.Columns().GatewayId, payment.GatewayId).OmitNil().Update()
	if err != nil {
		g.Log().Errorf(ctx, `SaveChannelUserDefaultPaymentMethod GatewayDefaultPaymentMethod failure %s`, err.Error())
	}
	return err
}

func HandlePaymentWebhookEvent(ctx context.Context, paymentId string, gatewayPaymentRo *gateway_bean.GatewayPaymentRo) error {
	one := query.GetPaymentByPaymentId(ctx, paymentId)
	if one == nil {
		one = query.GetPaymentByGatewayPaymentId(ctx, gatewayPaymentRo.GatewayPaymentId)
	}
	if one != nil {
		if gatewayPaymentRo.Status == consts.PaymentSuccess {
			err := HandlePaySuccess(ctx, &HandlePayReq{
				PaymentId:              one.PaymentId,
				GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
				GatewayUserId:          gatewayPaymentRo.GatewayUserId,
				TotalAmount:            gatewayPaymentRo.TotalAmount,
				PayStatusEnum:          consts.PaymentSuccess,
				PaidTime:               gatewayPaymentRo.PaidTime,
				PaymentAmount:          gatewayPaymentRo.PaymentAmount,
				CaptureAmount:          0,
				Reason:                 gatewayPaymentRo.Reason,
				GatewayPaymentMethod:   gatewayPaymentRo.GatewayPaymentMethod,
				PaymentCode:            gatewayPaymentRo.PaymentCode,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.Status == consts.PaymentFailed {
			err := HandlePayFailure(ctx, &HandlePayReq{
				PaymentId:              one.PaymentId,
				GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
				PayStatusEnum:          consts.PaymentFailed,
				Reason:                 gatewayPaymentRo.Reason,
				PaymentCode:            gatewayPaymentRo.PaymentCode,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.Status == consts.PaymentCancelled {
			err := HandlePayCancel(ctx, &HandlePayReq{
				PaymentId:              one.PaymentId,
				GatewayPaymentIntentId: gatewayPaymentRo.GatewayPaymentId,
				GatewayPaymentId:       gatewayPaymentRo.GatewayPaymentId,
				PayStatusEnum:          consts.PaymentCancelled,
				Reason:                 gatewayPaymentRo.Reason,
				PaymentCode:            gatewayPaymentRo.PaymentCode,
			})
			if err != nil {
				return err
			}
		} else if gatewayPaymentRo.AuthorizeStatus == consts.WaitingAuthorized {
			err := HandlePayNeedAuthorized(ctx, one, gatewayPaymentRo.AuthorizeReason, gatewayPaymentRo.PaymentData, gatewayPaymentRo.PaymentCode)
			if err != nil {
				return err
			}
		}
	} else {
		return gerror.Newf("Payment Not Match Or Not Found GatewayPaymentId:%s", gatewayPaymentRo.GatewayPaymentId)
	}

	return nil
}

func CreateOrUpdatePaymentTimelineForPayment(ctx context.Context, payment *entity.Payment, uniqueId string) error {
	one := query.GetPaymentTimeLineByUniqueId(ctx, uniqueId)
	payment = query.GetPaymentByPaymentId(ctx, payment.PaymentId)

	var status = 0
	if payment.Status == consts.PaymentSuccess {
		status = 1
	} else if payment.Status == consts.PaymentFailed || payment.Status == consts.PaymentCancelled {
		status = 2
	}
	if one == nil {
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
			TimelineType:   consts.TimelineTypePayment,
			CreateTime:     gtime.Now().Timestamp(),
		}
		result, err := dao.PaymentTimeline.Ctx(ctx).Data(one).OmitNil().Insert(one)
		if err != nil {
			err = gerror.Newf(`CreateOrUpdatePaymentTimelineForPayment record insert failure %s`, err.Error())
			return err
		}
		id, _ := result.LastInsertId()
		one.Id = uint64(id)
	} else {
		_, err := dao.PaymentTimeline.Ctx(ctx).Data(g.Map{
			dao.PaymentTimeline.Columns().GmtModify: gtime.Now(),
			dao.PaymentTimeline.Columns().Status:    status,
		}).Where(dao.PaymentTimeline.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			err = gerror.Newf(`CreateOrUpdatePaymentTimelineForPayment record update failure %s`, err.Error())
			return err
		}
	}
	invoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	if invoice != nil {
		err := CreateOrUpdatePaymentItemForPaymentInvoice(ctx, payment, invoice)
		if err != nil {
			g.Log().Errorf(ctx, "CreateOrUpdatePaymentItemForPaymentInvoice error:%s", err.Error())
		}
	}
	return nil
}

func CreateOrUpdatePaymentItemForPaymentInvoice(ctx context.Context, payment *entity.Payment, invoice *entity.Invoice) error {
	if payment == nil || invoice == nil {
		return gerror.Newf(`payment or invoice is nil`)
	}
	payment = query.GetPaymentByPaymentId(ctx, payment.PaymentId)
	var status = 0
	if payment.Status == consts.PaymentSuccess {
		status = 1
	} else if payment.Status == consts.PaymentFailed || payment.Status == consts.PaymentCancelled {
		status = 2
	}
	var list []*entity.PaymentItem
	var total = 0
	_ = dao.PaymentItem.Ctx(ctx).
		Where(dao.PaymentItem.Columns().PaymentId, payment).
		OmitEmpty().ScanAndCount(&list, &total, true)
	invoiceSimplify := bean.SimplifyInvoice(invoice)
	if total != len(invoiceSimplify.Lines) {
		_, _ = dao.PaymentItem.Ctx(ctx).Where(dao.PaymentItem.Columns().PaymentId, payment.PaymentId).Delete()
		for i, item := range invoiceSimplify.Lines {
			one := &entity.PaymentItem{
				BizType:        payment.BizType,
				MerchantId:     payment.MerchantId,
				UserId:         payment.UserId,
				SubscriptionId: payment.SubscriptionId,
				InvoiceId:      payment.InvoiceId,
				UniqueId:       fmt.Sprintf("%s_%d", payment.PaymentId, i),
				Currency:       payment.Currency,
				Name:           item.Name,
				Description:    item.Description,
				Amount:         item.Amount,
				UnitAmount:     item.UnitAmountExcludingTax,
				Quantity:       item.Quantity,
				GatewayId:      payment.GatewayId,
				PaymentId:      payment.PaymentId,
				Status:         status,
				CreateTime:     gtime.Now().Timestamp(),
			}
			_, err := dao.PaymentItem.Ctx(ctx).Data(one).OmitNil().Insert(one)
			if err != nil {
				return gerror.Newf(`record insert failure %s`, err.Error())
			}
		}
	} else {
		_, err := dao.PaymentItem.Ctx(ctx).Data(g.Map{
			dao.PaymentItem.Columns().GmtModify: gtime.Now(),
			dao.PaymentItem.Columns().Status:    status,
		}).Where(dao.PaymentItem.Columns().PaymentId, payment.PaymentId).OmitNil().Update()
		if err != nil {
			return gerror.Newf(`record update failure %s`, err.Error())
		}
	}

	return nil
}
