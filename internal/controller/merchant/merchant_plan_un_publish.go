package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) UnPublish(ctx context.Context, req *plan.UnPublishReq) (res *plan.UnPublishRes, err error) {
	err = service.SubscriptionPlanUnPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.UnPublishRes{}, nil
}
