package expire

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	handler2 "unibee-api/internal/logic/payment/handler"
	service2 "unibee-api/internal/logic/subscription/service"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/redismq"
)

func SubscriptionExpire(ctx context.Context, sub *entity.Subscription, reason string) error {
	//Expire SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		return err
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, reason)
		if err != nil {
			fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel Subscription Remaining Payment
	var paymentList []*entity.Payment
	err = dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().SubscriptionId, sub.SubscriptionId).
		Where(dao.Payment.Columns().Status, consts.TO_BE_PAID).
		Limit(0, 100).
		OmitEmpty().Scan(&paymentList)
	if err != nil {
		fmt.Printf("SubscriptionExpire GetPaymentList error:%s", err.Error())
	}
	for _, p := range paymentList {
		// todo mark should use PaymentGatewayCancel
		err := handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
			PaymentId:     p.PaymentId,
			PayStatusEnum: consts.PAY_CANCEL,
			Reason:        reason,
		})
		if err != nil {
			fmt.Printf("SubscriptionExpire HandlePayCancel error:%s", err.Error())
		}
	}
	//Expire Subscription UnFinished Invoice, May No Need
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:       consts.SubStatusExpired,
		dao.Subscription.Columns().CancelReason: reason,
		dao.Subscription.Columns().TrialEnd:     sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		fmt.Printf("SubscriptionExpire error:%s", err.Error())
		return err
	}

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionExpire.Topic,
		Tag:   redismq2.TopicSubscriptionExpire.Tag,
		Body:  sub.SubscriptionId,
	})

	return nil
}
