package addon

import (
	"context"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func GetSubscriptionAddonsByAddonJson(ctx context.Context, addonJson string) []*ro.PlanAddonVo {
	if len(addonJson) == 0 {
		return nil
	}
	var addonParams []*ro.SubscriptionPlanAddonParamRo
	err := utility.UnmarshalFromJsonString(addonJson, &addonParams)
	if err != nil {
		return nil
	}
	var addons []*ro.PlanAddonVo
	for _, param := range addonParams {
		addons = append(addons, &ro.PlanAddonVo{
			Quantity:  param.Quantity,
			AddonPlan: ro.SimplifyPlan(query.GetPlanById(ctx, param.AddonPlanId)),
		})
	}
	return addons
}