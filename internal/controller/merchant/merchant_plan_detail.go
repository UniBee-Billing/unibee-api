package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error) {

	return service.PlanDetail(ctx, _interface.GetMerchantId(ctx), req.PlanId)
}
