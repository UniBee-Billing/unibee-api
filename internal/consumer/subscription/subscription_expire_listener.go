package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	user2 "unibee/internal/consumer/webhook/user"
	dao "unibee/internal/dao/oversea_pay"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/subscription/handler"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/logic/subscription/user_sub_plan"
	"unibee/internal/logic/user"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionExpireListener struct {
}

func (t SubscriptionExpireListener) GetTopic() string {
	return redismq2.TopicSubscriptionExpire.Topic
}

func (t SubscriptionExpireListener) GetTag() string {
	return redismq2.TopicSubscriptionExpire.Tag
}

func (t SubscriptionExpireListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "SubscriptionExpireListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	if sub != nil {
		user.UpdateUserDefaultSubscriptionForUpdate(ctx, sub.UserId, sub.SubscriptionId)
	}
	//Cancelled SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		g.Log().Errorf(ctx, "SubscriptionCreatePaymentCheckListener Fetch PendingUpdateList Error:%s", err.Error())
		return redismq.ReconsumeLater
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, "SubscriptionExpire")
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionCreatePaymentCheckListener SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	//Cancel All Invoice
	handler2.CancelInvoiceForSubscription(ctx, sub)
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	handler.FinishOldTimelineBySubEnd(ctx, sub.SubscriptionId)
	subscription3.SendMerchantSubscriptionWebhookBackground(sub, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_EXPIRED)
	user2.SendMerchantUserMetricWebhookBackground(sub.UserId, event.UNIBEE_WEBHOOK_EVENT_USER_METRIC_UPDATED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionExpireListener())
	fmt.Println("SubscriptionExpireListener RegisterListener")
}

func NewSubscriptionExpireListener() *SubscriptionExpireListener {
	return &SubscriptionExpireListener{}
}
