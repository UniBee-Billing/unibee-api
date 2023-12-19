package query

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionPlanChannelBinding(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetActiveSubscriptionPlanChannelBinding(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId, Status: consts.PlanStatusActive}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
