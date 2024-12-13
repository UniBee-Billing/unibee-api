package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	"unibee/internal/logic/subscription/pending_update_cancel"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionUpdateListener struct {
}

func (t SubscriptionUpdateListener) GetTopic() string {
	return redismq2.TopicSubscriptionUpdate.Topic
}

func (t SubscriptionUpdateListener) GetTag() string {
	return redismq2.TopicSubscriptionUpdate.Tag
}

func (t SubscriptionUpdateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionUpdateListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
		//if len(sub.VatNumber) > 0 {
		//	sub_update.UpdateUserDefaultVatNumber(ctx, sub.UserId, sub.VatNumber)
		//}
		user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_UPDATED, message.CustomData)
		//user2.SendMerchantUserMetricWebhookBackground(sub.UserId, sub.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED, fmt.Sprintf("SubscriptionUpdate#%s", sub.SubscriptionId))
		if len(sub.PendingUpdateId) > 0 {
			err := pending_update_cancel.SubscriptionPendingUpdateCancel(ctx, sub.PendingUpdateId, "CancelByUpdate-"+sub.PendingUpdateId)
			if err != nil {
				g.Log().Errorf(ctx, "HandleSubscriptionNextBillingCycleUpdate SubscriptionPendingUpdateCancel pendingUpdateId:%s error:%s", sub.PendingUpdateId, err.Error())
			}
		}
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionUpdateListener())
	fmt.Println("SubscriptionUpdateListener RegisterListener")
}

func NewSubscriptionUpdateListener() *SubscriptionUpdateListener {
	return &SubscriptionUpdateListener{}
}
