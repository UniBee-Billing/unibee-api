package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Delete(ctx context.Context, req *plan.DeleteReq) (res *plan.DeleteRes, err error) {
	_, err = plan2.PlanDelete(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.DeleteRes{}, nil
}
