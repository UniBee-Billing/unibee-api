package merchant

import (
	"context"
	v1 "unibee/api/merchant/plan"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
)

func (c *ControllerPlan) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	plans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId:    _interface.BizCtx().Get(ctx).MerchantId,
		Type:          req.Type,
		Status:        req.Status,
		PublishStatus: req.PublishStatus,
		Currency:      req.Currency,
		SortField:     req.SortField,
		SortType:      req.SortType,
		Page:          req.Page,
		Count:         req.Count,
	})
	return &v1.ListRes{Plans: plans}, nil
}
