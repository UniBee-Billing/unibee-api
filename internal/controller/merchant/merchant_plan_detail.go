package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface/context"
	plan2 "unibee/internal/logic/plan"
)

func (c *ControllerPlan) Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error) {

	return plan2.PlanDetail(ctx, _interface.GetMerchantId(ctx), req.PlanId)
}
