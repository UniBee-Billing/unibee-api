package query

import (
	"context"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
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

func GetSubscriptionBySubscriptionId(ctx context.Context, subscriptionId string) (one *entity.Subscription) {
	err := dao.Subscription.Ctx(ctx).Where(entity.Subscription{SubscriptionId: subscriptionId}).OmitEmpty().Scan(&one)
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

func GetSubscriptionAddonsBySubscriptionId(ctx context.Context, subscriptionId string) []*ro.SubscriptionPlanAddonRo {
	one := GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	if one == nil || len(one.AddonData) == 0 {
		return nil
	}
	var addonParams []*ro.SubscriptionPlanAddonParamRo
	err := utility.UnmarshalFromJsonString(one.AddonData, &addonParams)
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

func GetSubscriptionPendingUpdatesBySubscriptionId(ctx context.Context, subscriptionId string) []*entity.SubscriptionPendingUpdate {
	var list []*entity.SubscriptionPendingUpdate
	err := dao.SubscriptionPendingUpdate.Ctx(ctx).Where(dao.SubscriptionPendingUpdate.Columns().SubscriptionId, subscriptionId).OmitEmpty().Scan(&list)
	if err != nil {
		return nil
	}
	return list

}
