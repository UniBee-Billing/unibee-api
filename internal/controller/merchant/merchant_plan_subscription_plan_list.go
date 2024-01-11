package merchant

import (
	"context"
	v1 "go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"
)

func (c *ControllerPlan) SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
	}
	
	plans := service.SubscriptionPlanList(ctx, &service.SubscriptionPlanListInternalReq{
		MerchantId: req.MerchantId,
		Type:       req.Type,
		Status:     req.Status,
		Currency:   req.Currency,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &v1.SubscriptionPlanListRes{Plans: plans}, nil
}
