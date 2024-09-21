package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Copy(ctx context.Context, req *plan.CopyReq) (res *plan.CopyRes, err error) {
	one, err := plan2.PlanCopy(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.CopyRes{Plan: bean.SimplifyPlan(one)}, nil
}
