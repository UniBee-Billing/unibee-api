package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	service2 "unibee-api/internal/logic/subscription/service"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
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
	redismq.RegisterListener(NewSubscriptionCancelListener())
	fmt.Println("SubscriptionCancelListener RegisterListener")
}

func NewSubscriptionCancelListener() *SubscriptionCancelListener {
	return &SubscriptionCancelListener{}
}
