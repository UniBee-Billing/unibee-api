package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface/context"
	plan2 "unibee/internal/logic/plan"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPlan) Copy(ctx context.Context, req *plan.CopyReq) (res *plan.CopyRes, err error) {
	if req.PlanId <= 0 {
		utility.Assert(len(req.ExternalPlanId) > 0, "either planId or externalPlanId should be set")
		one := query.GetPlanByExternalPlanId(ctx, _interface.GetMerchantId(ctx), req.ExternalPlanId)
		utility.Assert(one != nil, fmt.Sprintf("Plan not found by externalPlanId:%s", req.ExternalPlanId))
		req.PlanId = one.Id
	}
	one, err := plan2.PlanCopy(ctx, req.PlanId)
	if err != nil {
		return nil, err
	}
	return &plan.CopyRes{Plan: bean.SimplifyPlan(one)}, nil
}
