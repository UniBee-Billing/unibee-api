package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error) {

	return service.SubscriptionPlanDetail(ctx, req.PlanId)
}
