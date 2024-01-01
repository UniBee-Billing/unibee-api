package query

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetSubscriptionPlanChannel(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetActiveSubscriptionPlanChannel(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId, Status: consts.PlanStatusActive}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListActiveSubscriptionPlanChannels(ctx context.Context, planId int64) (list []*entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, Status: consts.PlanStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		list = nil
	}
	return
}
