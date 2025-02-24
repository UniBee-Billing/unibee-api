package event_charge

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func ComputeEventCharge(ctx context.Context, planId uint64, one *entity.MerchantMetricEvent) *bean.EventMetricCharge {
	plan := query.GetPlanById(ctx, planId)
	if plan == nil {
		return &bean.EventMetricCharge{}
	}
	met := query.GetMerchantMetric(ctx, one.MetricId)
	if met == nil {
		return &bean.EventMetricCharge{}
	}
	var chargingPrice *bean.MetricPlanChargeBindingParam
	list := bean.ConvertMetricPlanChargeListFromPlan(plan)
	for _, item := range list {
		if item.MetricId == one.MetricId {
			chargingPrice = item
		}
	}
	var amount int64 = 0
	if chargingPrice != nil && chargingPrice.ChargeType == 0 {
		amount = (int64(one.Used) - chargingPrice.StandardStartValue) * chargingPrice.StandardAmount
	} else if chargingPrice != nil && chargingPrice.ChargeType == 1 {
		// graduated pricing compute
	}

	return &bean.EventMetricCharge{
		PlanId:        plan.Id,
		CurrentValue:  one.Used,
		ChargePricing: chargingPrice,
		Amount:        amount,
		Currency:      plan.Currency,
	}
}
