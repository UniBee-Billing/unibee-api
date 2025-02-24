package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Archive(ctx context.Context, req *plan.ArchiveReq) (res *plan.ArchiveRes, err error) {
	_, err = plan2.PlanArchive(ctx, req.PlanId, req.HardArchive)
	if err != nil {
		return nil, err
	}
	return &plan.ArchiveRes{}, nil
}
