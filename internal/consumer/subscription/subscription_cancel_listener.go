package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
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
	"unibee/utility"
)

type SubscriptionCancelListener struct {
}

func (t SubscriptionCancelListener) GetTopic() string {
	return redismq2.TopicSubscriptionCancel.Topic
}

func (t SubscriptionCancelListener) GetTag() string {
	return redismq2.TopicSubscriptionCancel.Tag
}

func (t SubscriptionCancelListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionCancelListener Receive Message:%s", utility.MarshalToJsonString(message))
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
		g.Log().Errorf(ctx, "SubscriptionCancelListener Fetch SubscriptionPendingUpdate error:%s", err.Error())
		return redismq.ReconsumeLater
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.PendingUpdateId, "SubscriptionCancelled")
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionCancelListener SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel All Invoice
	service.TryCancelSubscriptionLatestInvoice(ctx, sub)
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	timeline.FinishOldTimelineBySubEnd(ctx, sub.SubscriptionId, consts.SubStatusCancelled)
	subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_CANCELLED, message.CustomData)
	user2.SendMerchantUserMetricWebhookBackground(sub.UserId, sub.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCancelListener())
	fmt.Println("SubscriptionCancelListener RegisterListener")
}

func NewSubscriptionCancelListener() *SubscriptionCancelListener {
	return &SubscriptionCancelListener{}
}
