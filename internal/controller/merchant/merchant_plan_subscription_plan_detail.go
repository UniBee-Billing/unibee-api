package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerPlan) SubscriptionPlanDetail(ctx context.Context, req *plan.SubscriptionPlanDetailReq) (res *plan.SubscriptionPlanDetailRes, err error) {
	return service.SubscriptionPlanDetail(ctx, req.PlanId)
}
