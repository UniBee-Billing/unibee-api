package event_charge

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"golang.org/x/text/currency"
)

func ComputeEventCharge(ctx context.Context, planId uint64, one *entity.MerchantMetricEvent, oldUsed int64, targetCurrency string) *bean.EventMetricCharge {
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
	oldTotalChargeAmount, unitAmount, graduatedStep, _ := ComputeMetricUsedChargePrice(ctx, plan, targetCurrency, oldUsed, chargingPrice)
	totalChargeAmount, unitAmount, graduatedStep, _ := ComputeMetricUsedChargePrice(ctx, plan, targetCurrency, one.Used, chargingPrice)
	return &bean.EventMetricCharge{
		PlanId:            plan.Id,
		CurrentUsedValue:  one.Used,
		ChargePricing:     chargingPrice,
		TotalChargeAmount: totalChargeAmount,
		ChargeAmount:      totalChargeAmount - oldTotalChargeAmount,
		UnitAmount:        unitAmount,
		GraduatedStep:     graduatedStep,
		Currency:          targetCurrency,
	}
}

func ComputeMetricUsedChargePrice(ctx context.Context, plan *entity.Plan, targetCurrency string, usedValue int64, chargingPrice *bean.PlanMetricMeteredChargeParam) (totalChargeAmount int64, unitAmount int64, graduatedStep *bean.MetricPlanChargeGraduatedStep, lone []*bean.MetricPlanChargeLine) {
	lines := make([]*bean.MetricPlanChargeLine, 0)
	var symbol = fmt.Sprintf("%v ", currency.NarrowSymbol(currency.MustParseISO(strings.ToUpper(targetCurrency))))
	if plan == nil {
		return 0, 0, nil, lines
	}
	totalChargeAmount = 0
	if chargingPrice != nil && chargingPrice.ChargeType == 0 && usedValue > 0 && chargingPrice.StandardAmount > 0 {
		unitAmount = plan.ExchangeAmountToCurrency(ctx, chargingPrice.StandardAmount, targetCurrency)
		totalChargeAmount = utility.MaxInt64(usedValue-chargingPrice.StandardStartValue, 0) * unitAmount
		lines = append(lines, &bean.MetricPlanChargeLine{
			UnitAmount: utility.MaxInt64(unitAmount, 0),
			Quantity:   utility.MaxInt64(usedValue-chargingPrice.StandardStartValue, 0),
			Amount:     utility.MaxInt64(totalChargeAmount, 0),
			FlatAmount: 0,
			Step:       "",
		})
	} else if chargingPrice != nil && chargingPrice.ChargeType == 1 && usedValue > 0 {
		var lastEnd int64 = 0
		for _, step := range chargingPrice.GraduatedAmounts {
			// reach end
			unitAmount = plan.ExchangeAmountToCurrency(ctx, step.PerAmount, targetCurrency)
			flatAmount := plan.ExchangeAmountToCurrency(ctx, step.FlatAmount, targetCurrency)
			if usedValue <= step.EndValue || step.EndValue < 0 {
				totalChargeAmount = ((usedValue - lastEnd) * unitAmount) + flatAmount + totalChargeAmount
				graduatedStep = step
				flatDesc := ""
				if flatAmount > 0 {
					flatDesc = fmt.Sprintf(" FlatAmount: %s%s ", symbol, utility.ConvertCentToDollarStr(flatAmount, targetCurrency))
				}
				lines = append(lines, &bean.MetricPlanChargeLine{
					UnitAmount: utility.MaxInt64(unitAmount, 0),
					Quantity:   utility.MaxInt64(usedValue-lastEnd, 0),
					Amount:     utility.MaxInt64(((usedValue-lastEnd)*unitAmount)+flatAmount, 0),
					FlatAmount: flatAmount,
					Step:       fmt.Sprintf("(%d - %d)%s", lastEnd, usedValue, flatDesc),
				})
				break
			} else {
				unitAmount = step.PerAmount
				totalChargeAmount = (step.EndValue-lastEnd)*unitAmount + flatAmount + totalChargeAmount
				graduatedStep = step
				flatDesc := ""
				if flatAmount > 0 {
					flatDesc = fmt.Sprintf(" FlatAmount: %s%s ", symbol, utility.ConvertCentToDollarStr(flatAmount, targetCurrency))
				}
				lines = append(lines, &bean.MetricPlanChargeLine{
					UnitAmount: utility.MaxInt64(unitAmount, 0),
					Quantity:   utility.MaxInt64(step.EndValue-lastEnd, 0),
					Amount:     utility.MaxInt64((step.EndValue-lastEnd)*unitAmount+flatAmount, 0),
					FlatAmount: flatAmount,
					Step:       fmt.Sprintf("(%d - %d)%s", lastEnd, step.EndValue, flatDesc),
				})
				lastEnd = step.EndValue
			}
		}
	}
	return utility.MaxInt64(totalChargeAmount, 0), utility.MaxInt64(unitAmount, 0), graduatedStep, lines
}
