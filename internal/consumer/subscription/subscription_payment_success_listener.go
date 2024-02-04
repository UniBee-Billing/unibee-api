package subscription

import (
	"context"
	"fmt"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/logic/user"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

type SubscriptionPaymentSuccessListener struct {
}

func (t SubscriptionPaymentSuccessListener) GetTopic() string {
	return redismq2.TopicSubscriptionPaymentSuccess.Topic
}

func (t SubscriptionPaymentSuccessListener) GetTag() string {
	return redismq2.TopicSubscriptionPaymentSuccess.Tag
}

func (t SubscriptionPaymentSuccessListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	fmt.Printf("SubscriptionPaymentSuccessListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscription(ctx, sub.UserId, sub.SubscriptionId)
		if len(sub.VatNumber) > 0 {
			user.UpdateUserDefaultVatNumber(ctx, sub.UserId, sub.VatNumber)
		}
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionPaymentSuccessListener())
	fmt.Println("SubscriptionPaymentSuccessListener RegisterListener")
}

func NewSubscriptionPaymentSuccessListener() *SubscriptionPaymentSuccessListener {
	return &SubscriptionPaymentSuccessListener{}
}
