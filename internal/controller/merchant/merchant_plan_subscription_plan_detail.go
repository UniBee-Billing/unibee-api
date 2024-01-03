package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/ro"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerPlan) SubscriptionPlanDetail(ctx context.Context, req *plan.SubscriptionPlanDetailReq) (res *plan.SubscriptionPlanDetailRes, err error) {
	one := query.GetSubscriptionPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, "plan not found")
	return &plan.SubscriptionPlanDetailRes{
		Plan: &ro.SubscriptionPlanRo{
			Plan:     one,
			Channels: query.GetListActiveSubscriptionPlanChannels(ctx, req.PlanId),
			Addons:   query.GetSubscriptionPlanAddonsByPlanId(ctx, req.PlanId),
		},
	}, nil
}
