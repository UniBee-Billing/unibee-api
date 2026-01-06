package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/onetime"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) OnetimeAddonPreview(ctx context.Context, req *subscription.OnetimeAddonPreviewReq) (res *subscription.OnetimeAddonPreviewRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		utility.Assert(req.AddonId > 0, "addonId should provide while SubscriptionId is blank")
		plan := query.GetPlanById(ctx, req.AddonId)
		utility.Assert(plan != nil, fmt.Sprintf("addon not found:%v", req.AddonId))
		sub := query.GetLatestSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), plan.ProductId)
		utility.Assert(sub != nil, "no subscription found")
		req.SubscriptionId = sub.SubscriptionId
	}
	preview, err := onetime.CreateSubOneTimeAddonPreview(ctx, &onetime.SubscriptionCreateOnetimeAddonInternalReq{
		MerchantId:             _interface.GetMerchantId(ctx),
		SubscriptionId:         req.SubscriptionId,
		AddonId:                req.AddonId,
		Currency:               req.Currency,
		Quantity:               req.Quantity,
		Metadata:               req.Metadata,
		DiscountCode:           req.DiscountCode,
		DiscountAmount:         req.DiscountAmount,
		DiscountPercentage:     req.DiscountPercentage,
		TaxPercentage:          req.TaxPercentage,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.OnetimeAddonPreviewRes{
		Addon:            bean.SimplifyPlan(preview.Addon),
		Quantity:         preview.Quantity,
		TaxAmount:        preview.Invoice.TaxAmount,
		DiscountAmount:   preview.Invoice.DiscountAmount,
		TotalAmount:      preview.Invoice.TotalAmount,
		OriginAmount:     preview.Invoice.OriginAmount,
		Currency:         preview.Invoice.Currency,
		Invoice:          preview.Invoice,
		UserId:           preview.User.Id,
		Email:            preview.User.Email,
		TaxPercentage:    preview.Invoice.TaxPercentage,
		VatNumber:        preview.Invoice.VatNumber,
		Discount:         preview.Invoice.Discount,
		ApplyPromoCredit: preview.ApplyPromoCredit,
	}, nil
}
