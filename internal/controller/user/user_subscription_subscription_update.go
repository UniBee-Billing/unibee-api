package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.EffectImmediate == 0, "EffectImmediate not support in user_portal")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)

	//Update 可以由 Admin 操作，service 层不做用户校验
	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id == sub.UserId, "userId not match")
	}
	return service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
		SubscriptionId:     req.SubscriptionId,
		NewPlanId:          req.NewPlanId,
		Quantity:           req.Quantity,
		GatewayId:          req.GatewayId,
		AddonParams:        req.AddonParams,
		ConfirmTotalAmount: req.ConfirmTotalAmount,
		ConfirmCurrency:    req.ConfirmCurrency,
		ProrationDate:      req.ProrationDate,
		EffectImmediate:    req.EffectImmediate,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
	}, 0)
}
