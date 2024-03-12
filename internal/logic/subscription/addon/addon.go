package addon

import (
	"context"
	"unibee/api/bean"
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
			AddonPlan: bean.SimplifyPlan(query.GetPlanById(ctx, param.AddonPlanId)),
		})
	}
	return addons
}
