package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/plan/service"
	"go-oversea-pay/utility"
)

func (c *ControllerPlan) SubscriptionPlanEdit(ctx context.Context, req *plan.SubscriptionPlanEditReq) (res *plan.SubscriptionPlanEditRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	one, err := service.SubscriptionPlanEdit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanEditRes{Plan: one}, nil
}
