package subscription_plan_merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/ro"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/subscription_plan_merchant/v1"
)

func (c *ControllerV1) SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error) {
	one := query.GetSubscriptionPlanById(ctx, req.PlanId)
	utility.Assert(one != nil, "plan not found")
	return &v1.SubscriptionPlanDetailRes{
		Plan: &ro.SubscriptionPlanRo{
			Plan:     one,
			Channels: query.GetListActiveSubscriptionPlanChannels(ctx, req.PlanId),
			Addons:   query.GetSubscriptionPlanAddonsByPlanId(ctx, req.PlanId),
		},
	}, nil
}
