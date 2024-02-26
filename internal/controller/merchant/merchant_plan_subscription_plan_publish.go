package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) SubscriptionPlanPublish(ctx context.Context, req *plan.SubscriptionPlanPublishReq) (res *plan.SubscriptionPlanPublishRes, err error) {
	err = service.SubscriptionPlanPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanPublishRes{}, nil
}
