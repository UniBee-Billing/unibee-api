package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	user2 "unibee/internal/consumer/webhook/user"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
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
	g.Log().Debugf(ctx, "SubscriptionPaymentSuccessListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscriptionForPaymentSuccess(ctx, sub.UserId, sub.SubscriptionId)
		if len(sub.VatNumber) > 0 {
			user.UpdateUserDefaultVatNumber(ctx, sub.UserId, sub.VatNumber)
		}
	}
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	subscription3.SendMerchantSubscriptionWebhookBackground(sub, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED)
	user2.SendMerchantUserMetricWebhookBackground(sub.UserId, event.MERCHANT_WEBHOOK_TAG_USER_METRIC_UPDATED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionPaymentSuccessListener())
	fmt.Println("SubscriptionPaymentSuccessListener RegisterListener")
}

func NewSubscriptionPaymentSuccessListener() *SubscriptionPaymentSuccessListener {
	return &SubscriptionPaymentSuccessListener{}
}
