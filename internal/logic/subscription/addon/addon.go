package addon

import (
	"context"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

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
			AddonPlan: query.GetPlanById(ctx, param.AddonPlanId),
		})
	}
	return addons
}
