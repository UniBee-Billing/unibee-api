package merchant

import (
	"context"
	v1 "unibee/api/merchant/plan"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/plan"
)

func (c *ControllerPlan) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	plans, total := plan.PlanList(ctx, &plan.ListInternalReq{
		MerchantId:    _interface.Context().Get(ctx).MerchantId,
		PlanIds:       req.PlanIds,
		ProductIds:    req.ProductIds,
		Type:          req.Type,
		Status:        req.Status,
		PublishStatus: req.PublishStatus,
		Currency:      req.Currency,
		SearchKey:     req.SearchKey,
		SortField:     req.SortField,
		SortType:      req.SortType,
		Page:          req.Page,
		Count:         req.Count,
	})
	return &v1.ListRes{Plans: plans, Total: total}, nil
}
