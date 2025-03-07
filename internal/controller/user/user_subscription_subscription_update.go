package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
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
	if !config.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id == sub.UserId, "userId not match")
	}
	resp, err := service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
		SubscriptionId:         req.SubscriptionId,
		NewPlanId:              req.NewPlanId,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		AddonParams:            req.AddonParams,
		ConfirmTotalAmount:     req.ConfirmTotalAmount,
		ConfirmCurrency:        req.ConfirmCurrency,
		ProrationDate:          req.ProrationDate,
		EffectImmediate:        req.EffectImmediate,
		Metadata:               req.Metadata,
		DiscountCode:           req.DiscountCode,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	}, 0)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: resp.SubscriptionPendingUpdate,
		Paid:                      resp.Paid,
		Link:                      resp.Link,
		Note:                      resp.Note,
	}, err
}
