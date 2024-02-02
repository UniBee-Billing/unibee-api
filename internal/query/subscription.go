package query

import (
	"context"
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

//func GetSubscriptionById(ctx context.Context, id int64) (one *entity.Subscription) {
//	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{Id: uint64(id)}).OmitEmpty().Scan(&one)
//	if err != nil {
//		one = nil
//	}
//	return
//}

func GetLatestActiveOrCreateSubscriptionByUserId(ctx context.Context, userId int64, merchantId int64) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).
		Where(entity.Subscription{UserId: userId}).
		Where(entity.Subscription{MerchantId: merchantId}).
		Where(entity.Subscription{IsDeleted: 0}).
		WhereIn(dao.Subscription.Columns().Status, []int{consts.SubStatusCreate, consts.SubStatusActive}).
		OrderDesc(dao.Subscription.Columns().GmtModify).
		OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionBySubscriptionId(ctx context.Context, subscriptionId string) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{SubscriptionId: subscriptionId}).Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionByChannelSubscriptionId(ctx context.Context, channelSubscriptionId string) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{ChannelSubscriptionId: channelSubscriptionId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetSubscriptionAddonsByAddonJson(ctx context.Context, addonJson string) []*ro.SubscriptionPlanAddonRo {
	if len(addonJson) == 0 {
		return nil
	}
	var addonParams []*ro.SubscriptionPlanAddonParamRo
	err := utility.UnmarshalFromJsonString(addonJson, &addonParams)
	if err != nil {
		return nil
	}
	var addons []*ro.SubscriptionPlanAddonRo
	for _, param := range addonParams {
		addons = append(addons, &ro.SubscriptionPlanAddonRo{
			Quantity:  param.Quantity,
			AddonPlan: GetPlanById(ctx, param.AddonPlanId),
		})
	}
	return addons
}

func GetSubscriptionUpgradePendingUpdateByChannelUpdateId(ctx context.Context, channelUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(channelUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().ChannelUpdateId, channelUpdateId).
		Where(dao.SubscriptionPendingUpdate.Columns().EffectImmediate, 1).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionDowngradePendingUpdateByPendingUpdateId(ctx context.Context, pendingUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		Where(dao.SubscriptionPendingUpdate.Columns().EffectImmediate, 0).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionPendingUpdateByPendingUpdateId(ctx context.Context, pendingUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUnfinishedSubscriptionPendingUpdateByPendingUpdateId(ctx context.Context, pendingUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(pendingUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().UpdateSubscriptionId, pendingUpdateId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUnfinishedSubscriptionPendingUpdateByChannelUpdateId(ctx context.Context, channelUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(channelUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().ChannelUpdateId, channelUpdateId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUnfinishedEffectImmediateSubscriptionPendingUpdateByChannelUpdateId(ctx context.Context, channelUpdateId string) *entity.SubscriptionPendingUpdate {
	if len(channelUpdateId) == 0 {
		return nil
	}
	var one *entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).
		Where(dao.SubscriptionPendingUpdate.Columns().ChannelUpdateId, channelUpdateId).
		WhereLT(dao.SubscriptionPendingUpdate.Columns().Status, consts.PendingSubStatusFinished).
		Where(dao.SubscriptionPendingUpdate.Columns().EffectImmediate, 1).
		OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetSubscriptionTimeLineByUniqueId(ctx context.Context, uniqueId string) (one *entity.SubscriptionTimeline) {
	err := dao.SubscriptionTimeline.Ctx(ctx).Where(entity.SubscriptionTimeline{UniqueId: uniqueId}).OmitEmpty().Scan(&one)
	if err != nil {
		one = nil
	}
	return
}
