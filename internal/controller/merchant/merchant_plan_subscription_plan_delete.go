package merchant

import (
	"context"
	"unibee/internal/logic/plan/service"

	"unibee/api/merchant/plan"
)

func (c *ControllerPlan) SubscriptionPlanDelete(ctx context.Context, req *plan.SubscriptionPlanDeleteReq) (res *plan.SubscriptionPlanDeleteRes, err error) {
	_, err = service.SubscriptionPlanDelete(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanDeleteRes{}, nil
}
