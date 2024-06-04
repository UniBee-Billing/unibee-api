package message

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/redismq"
	"unibee/utility"
)

type MerchantWebhookListener struct {
}

func (t MerchantWebhookListener) GetTopic() string {
	return redismq2.TopicMerchantWebhook.Topic
}

func (t MerchantWebhookListener) GetTag() string {
	return redismq2.TopicMerchantWebhook.Tag
}

func (t MerchantWebhookListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	if message.ReconsumeTimes > 7 {
		g.Log().Debugf(ctx, "Webhook_Subscription NewMerchantWebhookListener_Commit By Reach the Webhook limit 8 Message:%s", utility.MarshalToJsonString(message))
		return redismq.CommitMessage
	}
	g.Log().Debugf(ctx, "Webhook_Subscription NewMerchantWebhookListener Receive Message:%s", utility.MarshalToJsonString(message))
	var webhookMessage *WebhookMessage
	err := utility.UnmarshalFromJsonString(message.Body, &webhookMessage)

	if err != nil {
		g.Log().Errorf(ctx, "Webhook_Subscription NewMerchantWebhookListener_Resume By UnmarshalFromJsonString Error:%s", err.Error())
		return redismq.ReconsumeLater
	}

	if SendWebhookRequest(ctx, webhookMessage, message.ReconsumeTimes) {
		// todo mark limit the dependencyKey
		if webhookMessage.SequenceKey != "" {
			_, _ = g.Redis().Set(ctx, webhookMessage.SequenceKey, "Sent")
			_, _ = g.Redis().Expire(ctx, webhookMessage.SequenceKey, 24*60*60)
		}
		return redismq.CommitMessage
	}

	return redismq.ReconsumeLater
}

func init() {
	listener := NewMerchantWebhookListener()
	redismq.RegisterListener(listener)
}

func NewMerchantWebhookListener() *MerchantWebhookListener {
	return &MerchantWebhookListener{}
}
