package merchant

import (
	"context"
	v1 "unibee/api/merchant/plan"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	"unibee/utility"
)

func (c *ControllerPlan) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

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
