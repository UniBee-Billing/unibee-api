package detail

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/query"
)

type PlanDetail struct {
	Product               *bean.Product                   `json:"product" dc:"Product"`
	Plan                  *bean.Plan                      `json:"plan" dc:"Plan"`
	MetricPlanLimits      []*bean.MerchantMetricPlanLimit `json:"metricPlanLimits" dc:"MetricPlanLimits"`
	Addons                []*bean.Plan                    `json:"addons" dc:"Addons"`
	AddonIds              []int64                         `json:"addonIds" dc:"AddonIds"`
	OnetimeAddons         []*bean.Plan                    `json:"onetimeAddons" dc:"OneTimeAddons"`
	OnetimeAddonIds       []int64                         `json:"onetimeAddonIds" dc:"OneTimeAddonIds"`
	MetricMeteredCharge   []*MetricPlanChargeDetail       `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*MetricPlanChargeDetail       `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
}

type MetricPlanChargeDetail struct {
	MetricId           uint64                                `json:"metricId" dc:"MetricId" v:"required"`
	Metric             *bean.MerchantMetric                  `json:"merchantMetric"    description:"MerchantMetric"` // metricId
	ChargeType         int                                   `json:"chargeType" dc:"ChargeType,0-standard pricing 1-graduated pricing"`
	StandardAmount     int64                                 `json:"standardAmount" dc:"StandardAmount, used for standard pricing"`
	StandardStartValue int64                                 `json:"standardStartValue" dc:"StandardStartValue, used for standard pricing"`
	GraduatedAmounts   []*bean.MetricPlanChargeGraduatedStep `json:"graduatedAmounts" dc:"GraduatedAmounts"`
}

func ConvertMetricPlanChargeDetailArrayFromParam(ctx context.Context, params []*bean.MetricPlanChargeBindingParam) []*MetricPlanChargeDetail {
	var list = make([]*MetricPlanChargeDetail, 0)
	for _, param := range params {
		list = append(list, &MetricPlanChargeDetail{
			MetricId:           param.MetricId,
			Metric:             bean.SimplifyMerchantMetric(query.GetMerchantMetric(ctx, param.MetricId)),
			ChargeType:         param.ChargeType,
			StandardAmount:     param.StandardAmount,
			StandardStartValue: param.StandardStartValue,
			GraduatedAmounts:   param.GraduatedAmounts,
		})
	}
	return list
}
