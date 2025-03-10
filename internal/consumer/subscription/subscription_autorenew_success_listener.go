package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	"unibee/internal/query"
	"unibee/utility"
)

type SubscriptionAutoRenewSuccessListener struct {
}

func (t SubscriptionAutoRenewSuccessListener) GetTopic() string {
	return redismq2.TopicSubscriptionAutoRenewSuccess.Topic
}

func (t SubscriptionAutoRenewSuccessListener) GetTag() string {
	return redismq2.TopicSubscriptionAutoRenewSuccess.Tag
}

func (t SubscriptionAutoRenewSuccessListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionAutoRenewSuccessListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		sub = query.GetSubscriptionBySubscriptionId(ctx, sub.SubscriptionId)
		subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_AUTORENEW_SUCCESS, message.CustomData)
		//user2.SendMerchantUserMetricWebhookBackground(sub.UserId, sub.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED, fmt.Sprintf("SubscriptionAutoRenewSuccess#%s", sub.SubscriptionId))
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionAutoRenewSuccessListener())
	fmt.Println("SubscriptionAutoRenewSuccessListener RegisterListener")
}

func NewSubscriptionAutoRenewSuccessListener() *SubscriptionAutoRenewSuccessListener {
	return &SubscriptionAutoRenewSuccessListener{}
}
