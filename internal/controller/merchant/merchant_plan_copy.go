package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/logic/plan/service"

	"unibee/api/merchant/plan"
)

func (c *ControllerPlan) Copy(ctx context.Context, req *plan.CopyReq) (res *plan.CopyRes, err error) {
	one, err := service.PlanCopy(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.CopyRes{Plan: bean.SimplifyPlan(one)}, nil
}
