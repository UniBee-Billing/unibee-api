package merchant

import (
	"context"
	"unibee-api/api/merchant/plan"
	"unibee-api/internal/logic/plan/service"
)

func (c *ControllerPlan) SubscriptionPlanUnPublish(ctx context.Context, req *plan.SubscriptionPlanUnPublishReq) (res *plan.SubscriptionPlanUnPublishRes, err error) {
	err = service.SubscriptionPlanUnPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanUnPublishRes{}, nil
}
