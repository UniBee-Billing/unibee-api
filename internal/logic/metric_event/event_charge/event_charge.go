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
	usedAmount, unitAmount, graduatedStep := ComputeMetricUsedChargePrice(one.Used, chargingPrice)
	return &bean.EventMetricCharge{
		PlanId:        plan.Id,
		CurrentValue:  one.Used,
		ChargePricing: chargingPrice,
		UsedAmount:    usedAmount,
		UnitAmount:    unitAmount,
		GraduatedStep: graduatedStep,
		Currency:      plan.Currency,
	}
}

func ComputeMetricUsedChargePrice(usedValue uint64, chargingPrice *bean.MetricPlanChargeBindingParam) (usedAmount int64, unitAmount int64, graduatedStep *bean.MetricPlanChargeGraduatedStep) {
	usedAmount = 0
	if chargingPrice != nil && chargingPrice.ChargeType == 0 && usedValue > 0 && chargingPrice.StandardAmount > 0 {
		usedAmount = (int64(usedValue) - chargingPrice.StandardStartValue) * chargingPrice.StandardAmount
		unitAmount = chargingPrice.StandardAmount
	} else if chargingPrice != nil && chargingPrice.ChargeType == 1 && usedValue > 0 {
		var lastEnd int64 = 0
		for _, step := range chargingPrice.GraduatedAmounts {
			if int64(usedValue) <= step.EndValue || step.EndValue < 0 {
				// reach end
				usedAmount = ((int64(usedValue) - lastEnd) * step.PerAmount) + step.FlatAmount + usedAmount
				unitAmount = step.PerAmount
				graduatedStep = step
				break
			} else {
				usedAmount = (step.EndValue-lastEnd)*step.PerAmount + step.FlatAmount + usedAmount
				unitAmount = step.PerAmount
				graduatedStep = step
				lastEnd = step.EndValue
			}
		}
	}
	return usedAmount, unitAmount, graduatedStep
}
