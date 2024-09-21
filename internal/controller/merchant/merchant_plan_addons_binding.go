package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) AddonsBinding(ctx context.Context, req *plan.AddonsBindingReq) (res *plan.AddonsBindingRes, err error) {
	one, err := plan2.PlanAddonsBinding(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.AddonsBindingRes{Plan: bean.SimplifyPlan(one)}, nil
}
