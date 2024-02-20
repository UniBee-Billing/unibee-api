package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/logic/subscription/user_sub_plan"
	"unibee-api/internal/logic/user"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
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
	g.Log().Infof(ctx, "SubscriptionPaymentSuccessListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscription(ctx, sub.UserId, sub.SubscriptionId)
		if len(sub.VatNumber) > 0 {
			user.UpdateUserDefaultVatNumber(ctx, sub.UserId, sub.VatNumber)
		}
	}
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionPaymentSuccessListener())
	fmt.Println("SubscriptionPaymentSuccessListener RegisterListener")
}

func NewSubscriptionPaymentSuccessListener() *SubscriptionPaymentSuccessListener {
	return &SubscriptionPaymentSuccessListener{}
}
