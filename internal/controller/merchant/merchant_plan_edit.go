package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/plan"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	"unibee/utility"
)

func (c *ControllerPlan) Edit(ctx context.Context, req *plan.EditReq) (res *plan.EditRes, err error) {

	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	one, err := service.SubscriptionPlanEdit(ctx, req)
	if err != nil {
		return nil, err
	}
	return &plan.EditRes{Plan: bean.SimplifyPlan(one)}, nil
}
