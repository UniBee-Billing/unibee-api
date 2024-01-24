package query

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
)

func GetPlanChannel(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetActivePlanChannel(ctx context.Context, planId int64, channelId int64) (one *entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, ChannelId: channelId, Status: consts.PlanChannelStatusActive}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetListActivePlanChannels(ctx context.Context, planId int64) (list []*entity.SubscriptionPlanChannel) {
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, Status: consts.PlanChannelStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		list = nil
	}
	return
}

func GetListActiveOutChannelRos(ctx context.Context, planId int64) []*ro.OutChannelRo {
	var list []*entity.SubscriptionPlanChannel
	err := dao.SubscriptionPlanChannel.Ctx(ctx).Where(entity.SubscriptionPlanChannel{PlanId: planId, Status: consts.PlanChannelStatusActive}).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	var outChannels []*ro.OutChannelRo
	for _, planChannel := range list {
		if planChannel.Status == consts.PlanChannelStatusActive {
			outChannel := GetPayChannelById(ctx, planChannel.ChannelId)
			if outChannel != nil {
				outChannels = append(outChannels, &ro.OutChannelRo{
					ChannelId:   outChannel.Id,
					ChannelName: outChannel.Name,
				})
			}
		}
	}
	return outChannels
}
