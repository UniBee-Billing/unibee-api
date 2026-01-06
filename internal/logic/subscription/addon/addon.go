package addon

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func GetSubscriptionAddonsByAddonJson(ctx context.Context, addonJson string) []*bean.PlanAddonDetail {
	if len(addonJson) == 0 {
		return nil
	}
	var addonParams []*bean.PlanAddonParam
	err := utility.UnmarshalFromJsonString(addonJson, &addonParams)
	if err != nil {
		return nil
	}
	var addons []*bean.PlanAddonDetail
	for _, param := range addonParams {
		addons = append(addons, &bean.PlanAddonDetail{
			Quantity:  param.Quantity,
			AddonPlan: bean.SimplifyPlanWithContext(ctx, query.GetPlanById(ctx, param.AddonPlanId)),
		})
	}
	return addons
}

func GetSubscriptionOneTimeAddonsOfCurrentPeriod(ctx context.Context, sub *entity.Subscription) []*bean.PlanAddonDetail {
	list := make([]*bean.PlanAddonDetail, 0)
	oneTimeAddonPurchases := query.GetSubscriptionOnetimeAddons(ctx, sub)
	for _, oneTimeAddonPurchase := range oneTimeAddonPurchases {
		addon := query.GetPlanById(ctx, oneTimeAddonPurchase.AddonId)
		if addon != nil {
			list = append(list, &bean.PlanAddonDetail{
				Quantity:  oneTimeAddonPurchase.Quantity,
				AddonPlan: bean.SimplifyPlanWithContext(ctx, addon),
			})
		}
	}
	return list
}
