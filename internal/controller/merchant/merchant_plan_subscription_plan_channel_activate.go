package merchant

import (
	"context"
	"unibee/api/merchant/plan"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/plan/service"
	"unibee/utility"
)

func (c *ControllerPlan) SubscriptionPlanChannelActivate(ctx context.Context, req *plan.SubscriptionPlanChannelActivateReq) (res *plan.SubscriptionPlanChannelActivateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	err = service.SubscriptionGatewayPlanActivate(ctx, req.PlanId, req.GatewayId)
	if err != nil {
		return nil, err
	}
	return &plan.SubscriptionPlanChannelActivateRes{}, nil
}
