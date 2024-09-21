package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Publish(ctx context.Context, req *plan.PublishReq) (res *plan.PublishRes, err error) {
	err = plan2.PlanPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.PublishRes{}, nil
}
