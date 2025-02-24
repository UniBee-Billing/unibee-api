package bean

import (
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type MetricPlanChargeBindingParam struct {
	MetricId           uint64                           `json:"metricId" dc:"MetricId" v:"required"`
	ChargeType         int                              `json:"chargeType" dc:"ChargeType,0-standard pricing 1-graduated pricing"`
	StandardAmount     int64                            `json:"standardAmount" dc:"StandardAmount, used for standard pricing"`
	StandardStartValue int64                            `json:"standardStartValue" dc:"StandardStartValue, used for standard pricing"`
	GraduatedAmounts   []*MetricPlanChargeGraduatedStep `json:"graduatedAmounts" dc:"GraduatedAmounts, used for graduated pricing"`
}

type MetricPlanChargeGraduatedStep struct {
	PerAmount  float64 `json:"perAmount" dc:"PerAmount"`
	StartValue int64   `json:"startValue" dc:"StartValue"`
	EndValue   int64   `json:"endValue" dc:"EndValue"`
	FlatAmount float64 `json:"flatAmount" dc:"FlatAmount"`
}

type MetricPlanChargeEntity struct {
	MetricMeteredCharge   []*MetricPlanChargeBindingParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*MetricPlanChargeBindingParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
}

func ConvertMetricPlanChargeEntityFromPlan(one *entity.Plan) *MetricPlanChargeEntity {
	var metricPlanCharge = &MetricPlanChargeEntity{}
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &metricPlanCharge)
	}
	return metricPlanCharge
}

func ConvertMetricPlanChargeListFromPlan(one *entity.Plan) []*MetricPlanChargeBindingParam {
	var metricPlanCharge = &MetricPlanChargeEntity{}
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &metricPlanCharge)
	}
	var list = metricPlanCharge.MetricMeteredCharge
	for _, met := range metricPlanCharge.MetricMeteredCharge {
		list = append(list, met)
	}
	for _, met := range metricPlanCharge.MetricRecurringCharge {
		list = append(list, met)
	}
	return list
}

type EventMetricCharge struct {
	PlanId        uint64                        `json:"planId" dc:"PlanId"`
	CurrentValue  uint64                        `json:"currentValue" dc:"CurrentValue"`
	ChargePricing *MetricPlanChargeBindingParam `json:"chargePricing" dc:"ChargePricing"`
	Amount        int64                         `json:"amount" dc:"Amount"`
	Currency      string                        `json:"currency" dc:"Currency"`
}
