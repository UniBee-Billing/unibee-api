package detail

import "unibee/api/bean"

type PlanDetail struct {
	Product          *bean.Product                   `json:"product" dc:"Product"`
	Plan             *bean.Plan                      `json:"plan" dc:"Plan"`
	MetricPlanLimits []*bean.MerchantMetricPlanLimit `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	Addons           []*bean.Plan                    `json:"addons" dc:"Addons"`
	AddonIds         []int64                         `json:"addonIds" dc:"AddonIds"`
	OnetimeAddons    []*bean.Plan                    `json:"onetimeAddons" dc:"OneTimeAddons"`
	OnetimeAddonIds  []int64                         `json:"onetimeAddonIds" dc:"OneTimeAddonIds"`
}
