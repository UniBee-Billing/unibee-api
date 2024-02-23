package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consumer/webhook/event"
	subscription3 "unibee-api/internal/consumer/webhook/subscription"
	"unibee-api/internal/logic/subscription/user_sub_plan"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
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
	g.Log().Infof(ctx, "SubscriptionIncompleteListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	subscription3.SendSubscriptionMerchantWebhookBackground(sub, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionIncompleteListener())
	fmt.Println("SubscriptionIncompleteListener RegisterListener")
}

func NewSubscriptionIncompleteListener() *SubscriptionIncompleteListener {
	return &SubscriptionIncompleteListener{}
}
