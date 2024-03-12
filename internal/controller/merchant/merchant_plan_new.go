package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) New(ctx context.Context, req *plan.NewReq) (res *plan.NewRes, err error) {

	one, err := service.SubscriptionPlanCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.NewRes{Plan: bean.SimplifyPlan(one)}, nil
}
