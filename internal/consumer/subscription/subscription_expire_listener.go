package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consts"
	"unibee-api/internal/consumer/webhook/event"
	subscription3 "unibee-api/internal/consumer/webhook/subscription"
	dao "unibee-api/internal/dao/oversea_pay"
	service2 "unibee-api/internal/logic/subscription/service"
	"unibee-api/internal/logic/subscription/user_sub_plan"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
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
	g.Log().Infof(ctx, "SubscriptionExpireListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	//Cancelled SubscriptionPendingUpdate
	var pendingUpdates []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Limit(0, 100).
		OmitEmpty().Scan(&pendingUpdates)
	if err != nil {
		return redismq.ReconsumeLater
	}
	for _, p := range pendingUpdates {
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, "SubscriptionExpire")
		if err != nil {
			fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	user_sub_plan.ReloadUserSubPlanCacheListBackground(sub.MerchantId, sub.UserId)
	subscription3.SendSubscriptionMerchantWebhookBackground(sub, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionExpireListener())
	fmt.Println("SubscriptionExpireListener RegisterListener")
}

func NewSubscriptionExpireListener() *SubscriptionExpireListener {
	return &SubscriptionExpireListener{}
}