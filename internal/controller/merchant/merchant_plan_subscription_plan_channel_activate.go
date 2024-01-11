package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"
)

func (c *ControllerPlan) SubscriptionPlanChannelActivate(ctx context.Context, req *plan.SubscriptionPlanChannelActivateReq) (res *plan.SubscriptionPlanChannelActivateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
	}

	err = service.SubscriptionPlanChannelActivate(ctx, req.PlanId, req.ChannelId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanChannelActivateRes{}, nil
}
