package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/subscription/handler"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/logic/subscription/user_sub_plan"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type SubscriptionCancelListener struct {
}

func (t SubscriptionCancelListener) GetTopic() string {
	return redismq2.TopicSubscriptionCancel.Topic
}

func (t SubscriptionCancelListener) GetTag() string {
	return redismq2.TopicSubscriptionCancel.Tag
}

func (t SubscriptionCancelListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionCancelListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	//Cancelled SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		g.Log().Errorf(ctx, "SubscriptionCancelListener Fetch SubscriptionPendingUpdate error:%s", err.Error())
		return redismq.ReconsumeLater
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, "SubscriptionCancelled")
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionCancelListener SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	subscription3.SendSubscriptionMerchantWebhookBackground(sub, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED)
	handler.FinishOldTimelineBySubEnd(ctx, sub.SubscriptionId)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCancelListener())
	fmt.Println("SubscriptionCancelListener RegisterListener")
}

func NewSubscriptionCancelListener() *SubscriptionCancelListener {
	return &SubscriptionCancelListener{}
}
