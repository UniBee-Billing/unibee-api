package message

import (
	"context"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/utility"
)

type InternalWebhookListener struct {
}

func (t InternalWebhookListener) GetTopic() string {
	return redismq2.TopicInternalWebhook.Topic
}

func (t InternalWebhookListener) GetTag() string {
	return redismq2.TopicInternalWebhook.Tag
}

func (t InternalWebhookListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	return redismq.ReconsumeLater
}

func init() {
	listener := NewInternalWebhookListener()
	redismq.RegisterListener(listener)
}

func NewInternalWebhookListener() *InternalWebhookListener {
	return &InternalWebhookListener{}
}
