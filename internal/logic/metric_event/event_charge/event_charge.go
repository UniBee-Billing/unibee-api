package event_charge

import (
	"context"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func ComputeEventCharge(ctx context.Context, planId uint64, one *entity.MerchantMetricEvent, oldUsed uint64) *bean.EventMetricCharge {
	plan := query.GetPlanById(ctx, planId)
	if plan == nil {
		return &bean.EventMetricCharge{}
	}
	met := query.GetMerchantMetric(ctx, one.MetricId)
	if met == nil {
		return &bean.EventMetricCharge{}
	}
	var chargingPrice *bean.PlanMetricMeteredChargeParam
	list := bean.ConvertMetricPlanBindingListFromPlan(plan)
	for _, item := range list {
		if item.MetricId == one.MetricId {
			chargingPrice = item
		}
	}
	oldTotalChargeAmount, unitAmount, graduatedStep := ComputeMetricUsedChargePrice(oldUsed, chargingPrice)
	totalChargeAmount, unitAmount, graduatedStep := ComputeMetricUsedChargePrice(one.Used, chargingPrice)
	return &bean.EventMetricCharge{
		PlanId:            plan.Id,
		CurrentUsedValue:  one.Used,
		ChargePricing:     chargingPrice,
		TotalChargeAmount: totalChargeAmount,
		ChargeAmount:      totalChargeAmount - oldTotalChargeAmount,
		UnitAmount:        unitAmount,
		GraduatedStep:     graduatedStep,
		Currency:          plan.Currency,
	}
}

func ComputeMetricUsedChargePrice(usedValue uint64, chargingPrice *bean.PlanMetricMeteredChargeParam) (totalChargeAmount int64, unitAmount int64, graduatedStep *bean.MetricPlanChargeGraduatedStep) {
	totalChargeAmount = 0
	if chargingPrice != nil && chargingPrice.ChargeType == 0 && usedValue > 0 && chargingPrice.StandardAmount > 0 {
		totalChargeAmount = (int64(usedValue) - chargingPrice.StandardStartValue) * chargingPrice.StandardAmount
		unitAmount = chargingPrice.StandardAmount
	} else if chargingPrice != nil && chargingPrice.ChargeType == 1 && usedValue > 0 {
		var lastEnd int64 = 0
		for _, step := range chargingPrice.GraduatedAmounts {
			if int64(usedValue) <= step.EndValue || step.EndValue < 0 {
				// reach end
				totalChargeAmount = ((int64(usedValue) - lastEnd) * step.PerAmount) + step.FlatAmount + totalChargeAmount
				unitAmount = step.PerAmount
				graduatedStep = step
				break
			} else {
				totalChargeAmount = (step.EndValue-lastEnd)*step.PerAmount + step.FlatAmount + totalChargeAmount
				unitAmount = step.PerAmount
				graduatedStep = step
				lastEnd = step.EndValue
			}
		}
	}
	return totalChargeAmount, unitAmount, graduatedStep
}
