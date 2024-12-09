package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error) {
	update, err := service.SubscriptionUpdatePreview(ctx, &service.UpdatePreviewInternalReq{
		SubscriptionId:   req.SubscriptionId,
		NewPlanId:        req.NewPlanId,
		Quantity:         req.Quantity,
		AddonParams:      req.AddonParams,
		EffectImmediate:  req.EffectImmediate,
		DiscountCode:     req.DiscountCode,
		ApplyPromoCredit: req.ApplyPromoCredit,
	}, 0, -1)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdatePreviewRes{
		OriginAmount:      update.OriginAmount,
		TotalAmount:       update.TotalAmount,
		DiscountAmount:    update.DiscountAmount,
		Currency:          update.Currency,
		Invoice:           update.Invoice,
		NextPeriodInvoice: update.NextPeriodInvoice,
		ProrationDate:     update.ProrationDate,
		Discount:          update.Discount,
		DiscountMessage:   update.DiscountMessage,
		ApplyPromoCredit:  update.ApplyPromoCredit,
	}, nil
}
