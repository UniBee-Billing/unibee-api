package subscription

import (
	"context"
	"fmt"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
)

type SubscriptionCreateListener struct {
}

func (t SubscriptionCreateListener) GetTopic() string {
	return redismq2.TopicSubscriptionCancel.Topic
}

func (t SubscriptionCreateListener) GetTag() string {
	return redismq2.TopicSubscriptionCancel.Tag
}

func (t SubscriptionCreateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	fmt.Printf("SubscriptionCancelListener Receive Message:%s", utility.MarshalToJsonString(message))
	//sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)
	//Cancelled SubscriptionPendingUpdate
	//var pendingUpdates []*entity.SubscriptionPendingUpdate
	//err := dao.SubscriptionPendingUpdate.Ctx(ctx).
	//	Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, sub.SubscriptionId).
	//	WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
	//	Limit(0, 100).
	//	OmitEmpty().Scan(&pendingUpdates)
	//if err != nil {
	//	return redismq.ReconsumeLater
	//}
	//for _, p := range pendingUpdates {
	//	err = service2.SubscriptionPendingUpdateCancel(ctx, p.UpdateSubscriptionId, "SubscriptionCancelled")
	//	if err != nil {
	//		fmt.Printf("MakeSubscriptionExpired SubscriptionPendingUpdateCancel error:%s", err.Error())
	//	}
	//}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCreateListener())
	fmt.Println("SubscriptionCreateListener RegisterListener")
}

func NewSubscriptionCreateListener() *SubscriptionCreateListener {
	return &SubscriptionCreateListener{}
}
