package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerPlan) SubscriptionPlanChannelActivate(ctx context.Context, req *plan.SubscriptionPlanChannelActivateReq) (res *plan.SubscriptionPlanChannelActivateRes, err error) {
	err = service.SubscriptionPlanChannelActivate(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanChannelActivateRes{}, nil
}
