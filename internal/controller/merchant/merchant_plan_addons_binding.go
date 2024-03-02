package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/plan/service"
	"unibee/utility"
)

func (c *ControllerPlan) AddonsBinding(ctx context.Context, req *plan.AddonsBindingReq) (res *plan.AddonsBindingRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	one, err := service.SubscriptionPlanAddonsBinding(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.AddonsBindingRes{Plan: ro.SimplifyPlan(one)}, nil
}
