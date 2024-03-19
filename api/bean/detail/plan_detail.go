package detail

import "unibee/api/bean"

type PlanDetail struct {
	Plan             *bean.PlanSimplify              `json:"plan" dc:"Plan"`
	MetricPlanLimits []*bean.MerchantMetricPlanLimit `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	Addons           []*bean.PlanSimplify            `json:"addons" dc:"Addons"`
	AddonIds         []int64                         `json:"addonIds" dc:"AddonIds"`
	OnetimeAddons    []*bean.PlanSimplify            `json:"onetimeAddons" dc:"OneTimeAddons"`
	OnetimeAddonIds  []int64                         `json:"onetimeAddonIds" dc:"OneTimeAddonIds"`
}
