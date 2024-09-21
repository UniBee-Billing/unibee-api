package system

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/system/plan"
	plan2 "unibee/internal/logic/plan"
	"unibee/internal/query"
)

func (c *ControllerPlan) Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error) {
	one := query.GetPlanById(ctx, req.PlanId)
	if one != nil {
		detail, err := plan2.PlanDetail(ctx, one.MerchantId, one.Id)
		if err != nil {
			return nil, err
		}
		return &plan.DetailRes{Plan: detail.Plan}, nil
	}
	return nil, gerror.New("Plan Not Found")
}
