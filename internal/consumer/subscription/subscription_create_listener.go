package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
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
		sub_update.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
		user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
		_, _ = redismq.SendDelay(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionCreatePaymentCheck.Topic,
			Tag:        redismq2.TopicSubscriptionCreatePaymentCheck.Tag,
			Body:       sub.SubscriptionId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		}, 3*60)
		{
			sub.Status = consts.SubStatusPending
			subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_CREATED, message.CustomData)
			//user2.SendMerchantUserMetricWebhookBackground(sub.UserId, sub.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED, fmt.Sprintf("SubscriptionCreated#%s", sub.SubscriptionId))
			sub_update.UpdateUserVatNumber(ctx, sub.UserId, sub.VatNumber)
		}
		// 3min PaymentChecker
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCreateListener())
	fmt.Println("SubscriptionCreateListener RegisterListener")
}

func NewSubscriptionCreateListener() *SubscriptionCreateListener {
	return &SubscriptionCreateListener{}
}
