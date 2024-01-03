package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerPlan) SubscriptionPlanChannelDeactivate(ctx context.Context, req *plan.SubscriptionPlanChannelDeactivateReq) (res *plan.SubscriptionPlanChannelDeactivateRes, err error) {
	err = service.SubscriptionPlanChannelDeactivate(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanChannelDeactivateRes{}, nil
}
