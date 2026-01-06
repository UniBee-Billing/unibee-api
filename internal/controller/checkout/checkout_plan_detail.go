package checkout

import (
	"context"
	"unibee/api/bean"
	plan2 "unibee/internal/logic/plan"
	"unibee/internal/query"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/checkout/plan"
)

func (c *ControllerPlan) Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error) {
	one := query.GetPlanById(ctx, req.PlanId)
	if one != nil {
		detail, err := plan2.PlanDetail(ctx, one.MerchantId, one.Id)
		if err != nil {
			return nil, err
		}
		{
			var filterMultiCurrencies = make([]*bean.PlanMultiCurrency, 0)
			for _, multiCurrency := range detail.Plan.Plan.MultiCurrencies {
				if !multiCurrency.Disable {
					filterMultiCurrencies = append(filterMultiCurrencies, multiCurrency)
				}
			}
			detail.Plan.Plan.MultiCurrencies = filterMultiCurrencies
		}
		return &plan.DetailRes{Plan: detail.Plan, Merchant: bean.SimplifyMerchant(query.GetMerchantById(ctx, detail.Plan.Plan.MerchantId))}, nil
	}
	return nil, gerror.New("Plan Not Found")
}
