package merchant

import (
	"context"
	"unibee/internal/logic/plan/service"

	"unibee/api/merchant/plan"
)

func (c *ControllerPlan) Delete(ctx context.Context, req *plan.DeleteReq) (res *plan.DeleteRes, err error) {
	_, err = service.SubscriptionPlanDelete(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.DeleteRes{}, nil
}
