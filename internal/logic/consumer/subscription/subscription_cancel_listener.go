package subscription

import (
	"context"
	"fmt"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	service2 "go-oversea-pay/internal/logic/subscription/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
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
	fmt.Printf("SubscriptionCancelListener Receive Message:%s", utility.MarshalToJsonString(message))
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
		err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, "SubscriptionCancelled")
		if err != nil {
			fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
		}
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(New())
	fmt.Println("TestMessageListener RegisterListener")
}

func New() *SubscriptionCancelListener {
	return &SubscriptionCancelListener{}
}
