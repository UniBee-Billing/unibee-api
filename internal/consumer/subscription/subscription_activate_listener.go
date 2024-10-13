package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	user2 "unibee/internal/consumer/webhook/user"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionActiveListener struct {
}

func (t SubscriptionActiveListener) GetTopic() string {
	return redismq2.TopicSubscriptionActive.Topic
}

func (t SubscriptionActiveListener) GetTag() string {
	return redismq2.TopicSubscriptionActive.Tag
}

func (t SubscriptionActiveListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionActiveListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		sub_update.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, sub.UserId, sub.SubscriptionId)
		user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ACTIVATED, message.CustomData)
		user2.SendMerchantUserMetricWebhookBackground(sub.UserId, sub.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionActiveListener())
	fmt.Println("SubscriptionActiveListener RegisterListener")
}

func NewSubscriptionActiveListener() *SubscriptionActiveListener {
	return &SubscriptionActiveListener{}
}
