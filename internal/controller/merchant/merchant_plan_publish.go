package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) Publish(ctx context.Context, req *plan.PublishReq) (res *plan.PublishRes, err error) {
	err = service.SubscriptionPlanPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.PublishRes{}, nil
}
