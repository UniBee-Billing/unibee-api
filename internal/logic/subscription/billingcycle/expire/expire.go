package expire

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/operation_log"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	service2 "unibee/internal/logic/subscription/pending_update_cancel"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func SubscriptionExpire(ctx context.Context, sub *entity.Subscription, reason string) error {
	sub = query.GetSubscriptionBySubscriptionId(ctx, sub.SubscriptionId)
	if sub == nil {
		return gerror.New("sub not found")
	}
	if sub.Status == consts.SubStatusCancelled {
		return gerror.New("sub already cancelled")
	}
	if sub.Status == consts.SubStatusExpired {
		return gerror.New("sub already expired")
	}
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
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.PendingUpdateId, reason)
		if err != nil {
			fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel Subscription Remaining Payment
	var paymentList []*entity.Payment
	err = dao.Payment.Ctx(ctx).
		Where(dao.Payment.Columns().SubscriptionId, sub.SubscriptionId).
		Where(dao.Payment.Columns().Status, consts.PaymentCreated).
		Limit(0, 100).
		OmitEmpty().Scan(&paymentList)
	if err != nil {
		fmt.Printf("SubscriptionExpire GetPaymentList error:%s", err.Error())
	}
	for _, p := range paymentList {
		err = service.PaymentGatewayCancel(ctx, p)
		if err != nil {
			fmt.Printf("SubscriptionExpire PaymentGatewayCancel error:%s", err.Error())
		}
		err = handler2.HandlePayCancel(ctx, &handler2.HandlePayReq{
			PaymentId:     p.PaymentId,
			PayStatusEnum: consts.PaymentCancelled,
			Reason:        reason,
		})
		if err != nil {
			fmt.Printf("SubscriptionExpire HandlePayCancel error:%s", err.Error())
		}
	}
	nextStatus := consts.SubStatusExpired
	if sub.FirstPaidTime == 0 {
		nextStatus = consts.SubStatusFailed
	}
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:         nextStatus,
		dao.Subscription.Columns().CancelReason:   reason,
		dao.Subscription.Columns().TrialEnd:       sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", sub.SubscriptionId),
		Content:        fmt.Sprintf("%s(%s->%s)", reason, consts.SubStatusToEnum(sub.Status).Description(), consts.SubStatusToEnum(nextStatus).Description()),
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		fmt.Printf("SubscriptionExpireOrFailed error:%s", err.Error())
		return err
	}

	if nextStatus == consts.SubStatusExpired {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionExpire.Topic,
			Tag:        redismq2.TopicSubscriptionExpire.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	} else {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionFailed.Topic,
			Tag:        redismq2.TopicSubscriptionFailed.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}

	return nil
}
