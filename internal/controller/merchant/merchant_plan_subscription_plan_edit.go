package merchant

import (
	"context"
	"unibee-api/api/merchant/plan"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/plan/service"
	"unibee-api/utility"
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
