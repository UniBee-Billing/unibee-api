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

func (c *ControllerSubscription) UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.EffectImmediate == 0, "EffectImmediate not support in user_portal")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	if !config.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id == sub.UserId, "userId not match")
	}
	prepare, err := service.SubscriptionUpdatePreview(ctx, &service.UpdatePreviewInternalReq{
		SubscriptionId:         req.SubscriptionId,
		NewPlanId:              req.NewPlanId,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		EffectImmediate:        req.EffectImmediate,
		AddonParams:            req.AddonParams,
		DiscountCode:           req.DiscountCode,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	}, 0, 0)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdatePreviewRes{
		OriginAmount:      prepare.OriginAmount,
		TotalAmount:       prepare.TotalAmount,
		DiscountAmount:    prepare.DiscountAmount,
		Currency:          prepare.Currency,
		Invoice:           prepare.Invoice,
		NextPeriodInvoice: prepare.NextPeriodInvoice,
		ProrationDate:     prepare.ProrationDate,
		Discount:          prepare.Discount,
		DiscountMessage:   prepare.DiscountMessage,
		ApplyPromoCredit:  prepare.ApplyPromoCredit,
	}, nil
}
