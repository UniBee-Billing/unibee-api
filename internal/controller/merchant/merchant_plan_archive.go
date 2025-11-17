package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface/context"
	plan2 "unibee/internal/logic/plan"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPlan) Archive(ctx context.Context, req *plan.ArchiveReq) (res *plan.ArchiveRes, err error) {
	if req.PlanId <= 0 {
		utility.Assert(len(req.ExternalPlanId) > 0, "either planId or externalPlanId should be set")
		one := query.GetPlanByExternalPlanId(ctx, _interface.GetMerchantId(ctx), req.ExternalPlanId)
		utility.Assert(one != nil, fmt.Sprintf("Plan not found by externalPlanId:%s", req.ExternalPlanId))
		req.PlanId = one.Id
	}
	_, err = plan2.PlanArchive(ctx, req.PlanId, req.HardArchive)
	if err != nil {
		return nil, err
	}
	return &plan.ArchiveRes{}, nil
}
