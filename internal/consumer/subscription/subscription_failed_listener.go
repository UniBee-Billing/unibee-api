package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	user2 "unibee/internal/consumer/webhook/user"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/invoice/service"
	service2 "unibee/internal/logic/subscription/pending_update_cancel"
	"unibee/internal/logic/subscription/timeline"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user/sub_update"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionFailedListener struct {
}

func (t SubscriptionFailedListener) GetTopic() string {
	return redismq2.TopicSubscriptionFailed.Topic
}

func (t SubscriptionFailedListener) GetTag() string {
	return redismq2.TopicSubscriptionFailed.Tag
}

func (t SubscriptionFailedListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "SubscriptionFailedListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
	}
	//Cancelled SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		g.Log().Errorf(ctx, "SubscriptionCreatePaymentCheckListener Fetch PendingUpdateList Error:%s", err.Error())
		return redismq.ReconsumeLater
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.PendingUpdateId, "SubscriptionFailed")
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionCreatePaymentCheckListener SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel All Invoice
	service.TryCancelSubscriptionLatestInvoice(ctx, sub)
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	timeline.FinishOldTimelineBySubEnd(ctx, sub.SubscriptionId, consts.SubStatusFailed)
	subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_FAILED)
	user2.SendMerchantUserMetricWebhookBackground(sub.UserId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionFailedListener())
	fmt.Println("SubscriptionFailedListener RegisterListener")
}

func NewSubscriptionFailedListener() *SubscriptionFailedListener {
	return &SubscriptionFailedListener{}
}
