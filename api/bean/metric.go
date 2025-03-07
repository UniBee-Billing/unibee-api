package bean

import (
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

type PlanMetricLimitParam struct {
	MetricId    uint64 `json:"metricId" dc:"MetricId" v:"required"`
	MetricLimit uint64 `json:"metricLimit" dc:"MetricLimit" v:"required"`
}

type PlanMetricMeteredChargeParam struct {
	MetricId           uint64                           `json:"metricId" dc:"MetricId"`
	ChargeType         int                              `json:"chargeType" dc:"ChargeType,0-standard pricing 1-graduated pricing"`
	StandardAmount     int64                            `json:"standardAmount" dc:"StandardAmount, used for standard pricing,cent"`
	StandardStartValue int64                            `json:"standardStartValue" dc:"StandardStartValue, used for standard pricing"`
	GraduatedAmounts   []*MetricPlanChargeGraduatedStep `json:"graduatedAmounts" dc:"GraduatedAmounts, used for graduated pricing"`
}

type MetricPlanChargeGraduatedStep struct {
	PerAmount  int64 `json:"perAmount" dc:"PerAmount,cent"`
	StartValue int64 `json:"startValue" dc:"StartValue"`
	EndValue   int64 `json:"endValue" dc:"EndValue, -1 = infinity value(∞)"`
	FlatAmount int64 `json:"flatAmount" dc:"FlatAmount,cent"`
}

type MetricPlanBindingEntity struct {
	MetricLimits          []*PlanMetricLimitParam         `json:"metricLimits"  dc:"Plan's MetricLimit List" `
	MetricMeteredCharge   []*PlanMetricMeteredChargeParam `json:"metricMeteredCharge"  dc:"Plan's MetricMeteredCharge" `
	MetricRecurringCharge []*PlanMetricMeteredChargeParam `json:"metricRecurringCharge"  dc:"Plan's MetricRecurringCharge" `
}

func ConvertMetricPlanBindingEntityFromPlan(one *entity.Plan) *MetricPlanBindingEntity {
	var metricPlanCharge = &MetricPlanBindingEntity{}
	if len(one.MetricCharge) > 0 {
		_ = utility.UnmarshalFromJsonString(one.MetricCharge, &metricPlanCharge)
	}
	return metricPlanCharge
}

func ConvertMetricPlanBindingListFromPlan(one *entity.Plan) []*PlanMetricMeteredChargeParam {
	var metricPlanCharge = &MetricPlanBindingEntity{}
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
	PlanId            uint64                         `json:"planId" dc:"PlanId"`
	CurrentUsedValue  int64                          `json:"currentValue" dc:"CurrentUsedValue"`
	ChargePricing     *PlanMetricMeteredChargeParam  `json:"chargePricing" dc:"ChargePricing"`
	TotalChargeAmount int64                          `json:"totalChargeAmount" dc:"TotalChargeAmount"`
	ChargeAmount      int64                          `json:"chargeAmount" dc:"ChargeAmount"`
	UnitAmount        int64                          `json:"unitAmount" dc:"UnitAmount"`
	GraduatedStep     *MetricPlanChargeGraduatedStep `json:"graduatedStep" dc:"GraduatedStep"`
	Currency          string                         `json:"currency" dc:"Currency"`
}

type UserMetricChargeInvoiceItemEntity struct {
	MeteredChargeStats   []*UserMetricChargeInvoiceItem `json:"meteredChargeStats" dc:"MeteredChargeStats"`
	RecurringChargeStats []*UserMetricChargeInvoiceItem `json:"recurringChargeStats" dc:"RecurringChargeStats"`
}
type UserMetricChargeInvoiceItem struct {
	MetricId          uint64                        `json:"metricId" dc:"MetricId" v:"required"`
	CurrentUsedValue  int64                         `json:"CurrentUsedValue" dc:"CurrentUsedValue"`
	MaxEventId        uint64                        `json:"maxEventId"`
	MinEventId        uint64                        `json:"minEventId"`
	ChargePricing     *PlanMetricMeteredChargeParam `json:"chargePricing" dc:"ChargePricing"`
	TotalChargeAmount int64                         `json:"totalChargeAmount" dc:"TotalChargeAmount"`
	Name              string                        `json:"name" dc:"Name"`
	Description       string                        `json:"description" dc:"Description"`
}
