package subscription

import (
	"context"
	"fmt"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
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
	fmt.Printf("SubscriptionCreateListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	_, _ = redismq.SendDelay(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCreate.Topic,
		Tag:   redismq2.TopicSubscriptionCreate.Tag,
		Body:  sub.SubscriptionId,
	}, 3*60)
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
