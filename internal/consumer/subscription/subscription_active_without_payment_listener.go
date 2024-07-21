package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	user2 "unibee/internal/consumer/webhook/user"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionActiveWithoutPaymentListener struct {
}

func (t SubscriptionActiveWithoutPaymentListener) GetTopic() string {
	return redismq2.TopicSubscriptionActiveWithoutPayment.Topic
}

func (t SubscriptionActiveWithoutPaymentListener) GetTag() string {
	return redismq2.TopicSubscriptionActiveWithoutPayment.Tag
}

func (t SubscriptionActiveWithoutPaymentListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "SubscriptionActiveWithoutPaymentListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
		user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_UPDATED)
		user2.SendMerchantUserMetricWebhookBackground(sub.UserId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionActiveWithoutPaymentListener())
	fmt.Println("SubscriptionActiveWithoutPaymentListener RegisterListener")
}

func NewSubscriptionActiveWithoutPaymentListener() *SubscriptionActiveWithoutPaymentListener {
	return &SubscriptionActiveWithoutPaymentListener{}
}
