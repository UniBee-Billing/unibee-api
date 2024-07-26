package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) UnPublish(ctx context.Context, req *plan.UnPublishReq) (res *plan.UnPublishRes, err error) {
	err = plan2.PlanUnPublish(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.UnPublishRes{}, nil
}
