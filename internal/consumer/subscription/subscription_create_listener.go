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
	"unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionCreateListener struct {
}

func (t SubscriptionCreateListener) GetTopic() string {
	return redismq2.TopicSubscriptionCreate.Topic
}

func (t SubscriptionCreateListener) GetTag() string {
	return redismq2.TopicSubscriptionCreate.Tag
}

func (t SubscriptionCreateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "SubscriptionCreateListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
	}
	_, _ = redismq.SendDelay(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCreatePaymentCheck.Topic,
		Tag:   redismq2.TopicSubscriptionCreatePaymentCheck.Tag,
		Body:  sub.SubscriptionId,
	}, 3*60)

	{
		sub.Status = consts.SubStatusCreate
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED)
		user2.SendMerchantUserMetricWebhookBackground(sub.UserId, event.MERCHANT_WEBHOOK_TAG_USER_METRIC_UPDATED)
	}

	// 3min PaymentChecker
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCreateListener())
	fmt.Println("SubscriptionCreateListener RegisterListener")
}

func NewSubscriptionCreateListener() *SubscriptionCreateListener {
	return &SubscriptionCreateListener{}
}
