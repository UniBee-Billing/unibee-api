package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) AddonsBinding(ctx context.Context, req *plan.AddonsBindingReq) (res *plan.AddonsBindingRes, err error) {
	one, err := service.SubscriptionPlanAddonsBinding(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.AddonsBindingRes{Plan: bean.SimplifyPlan(one)}, nil
}
