package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionIncompleteListener struct {
}

func (t SubscriptionIncompleteListener) GetTopic() string {
	return redismq2.TopicSubscriptionIncomplete.Topic
}

func (t SubscriptionIncompleteListener) GetTag() string {
	return redismq2.TopicSubscriptionIncomplete.Tag
}

func (t SubscriptionIncompleteListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "SubscriptionIncompleteListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
		user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_UPDATED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionIncompleteListener())
	fmt.Println("SubscriptionIncompleteListener RegisterListener")
}

func NewSubscriptionIncompleteListener() *SubscriptionIncompleteListener {
	return &SubscriptionIncompleteListener{}
}
