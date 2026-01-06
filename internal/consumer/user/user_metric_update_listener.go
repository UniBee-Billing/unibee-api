package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consumer/webhook/event"
	user2 "unibee/internal/consumer/webhook/user"
	"unibee/internal/logic/metric"
	"unibee/utility"
)

type UserMetricUpdateListener struct {
}

func (t UserMetricUpdateListener) GetTopic() string {
	return redismq2.TopicUserMetricUpdate.Topic
}

func (t UserMetricUpdateListener) GetTag() string {
	return redismq2.TopicUserMetricUpdate.Tag
}

func (t UserMetricUpdateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "UserMetricUpdateListener Receive Message:%s", utility.MarshalToJsonString(message))
	if len(message.Body) > 0 {
		var one = &metric.UserMetricUpdateMessage{}
		err := utility.UnmarshalFromJsonString(message.Body, &one)
		if err != nil {
			g.Log().Errorf(ctx, "UserMetricUpdateListener Receive Message:%s UnmarshalFromJsonString error:%s", utility.MarshalToJsonString(message), err.Error())
			return redismq.CommitMessage
		}
		if one == nil || one.UserId <= 0 || len(one.SubscriptionId) == 0 {
			g.Log().Errorf(ctx, "UserMetricUpdateListener Receive Message:%s invalid data one:%s", utility.MarshalToJsonString(message), utility.MarshalToJsonString(one))
			return redismq.CommitMessage
		}
		user2.SendMerchantUserMetricWebhookBackground(one.UserId, one.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED, fmt.Sprintf("%s#%s", one.Description, one.SubscriptionId))
		user2.SendMerchantUserMetricWebhookBackground(one.UserId, one.SubscriptionId, event.UNIBEE_WEBHOOK_EVENT_USER_SUBSCRIPTION_UPDATED, fmt.Sprintf("%s#%s", one.Description, one.SubscriptionId))
	}

	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewUserMetricUpdateListener())
	fmt.Println("UserMetricUpdateListener RegisterListener")
}

func NewUserMetricUpdateListener() *UserMetricUpdateListener {
	return &UserMetricUpdateListener{}
}
