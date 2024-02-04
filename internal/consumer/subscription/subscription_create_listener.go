package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
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
	g.Log().Infof(ctx, "SubscriptionCreateListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	_, _ = redismq.SendDelay(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCreatePaymentCheck.Topic,
		Tag:   redismq2.TopicSubscriptionCreatePaymentCheck.Tag,
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
