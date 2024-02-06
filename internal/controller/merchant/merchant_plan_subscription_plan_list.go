package merchant

import (
	"context"
	v1 "unibee-api/api/merchant/plan"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/plan/service"
	"unibee-api/utility"
)

func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	plans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId:    req.MerchantId,
		Type:          req.Type,
		Status:        req.Status,
		PublishStatus: req.PublishStatus,
		Currency:      req.Currency,
		SortField:     req.SortField,
		SortType:      req.SortType,
		Page:          req.Page,
		Count:         req.Count,
	})
	return &v1.SubscriptionPlanListRes{Plans: plans}, nil
}
