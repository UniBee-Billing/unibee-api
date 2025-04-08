package checkout

import (
	"context"
	subscription2 "unibee/api/merchant/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/checkout/subscription"
	"unibee/internal/controller/merchant"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	utility.Assert(req.PlanId > 0, "PlanId is required")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "Plan not found")
	_interface.Context().Get(ctx).MerchantId = plan.MerchantId
	controllerSubscription := merchant.ControllerSubscription{}
	preview, err := controllerSubscription.CreatePreview(ctx, &subscription2.CreatePreviewReq{
		PlanId:                 req.PlanId,
		Email:                  req.Email,
		UserId:                 req.UserId,
		ExternalUserId:         req.ExternalUserId,
		User:                   req.User,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		AddonParams:            req.AddonParams,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		TaxPercentage:          req.TaxPercentage,
		DiscountCode:           req.DiscountCode,
		TrialEnd:               req.TrialEnd,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.CreatePreviewRes{
		Plan:                           preview.Plan,
		TrialEnd:                       preview.TrialEnd,
		Quantity:                       preview.Quantity,
		Gateway:                        preview.Gateway,
		AddonParams:                    preview.AddonParams,
		Addons:                         preview.Addons,
		TaxPercentage:                  preview.TaxPercentage,
		SubscriptionAmountExcludingTax: preview.Invoice.SubscriptionAmountExcludingTax,
		TaxAmount:                      preview.Invoice.TaxAmount,
		DiscountAmount:                 preview.DiscountAmount,
		TotalAmount:                    preview.TotalAmount,
		OriginAmount:                   preview.OriginAmount,
		Currency:                       preview.Currency,
		VatNumber:                      preview.VatNumber,
		VatNumberValidate:              preview.VatNumberValidate,
		VatCountryCode:                 preview.VatCountryCode,
		VatCountryName:                 preview.VatCountryName,
		Invoice:                        preview.Invoice,
		UserId:                         preview.UserId,
		Email:                          preview.Email,
		Discount:                       preview.Discount,
		VatNumberValidateMessage:       preview.VatNumberValidateMessage,
		DiscountMessage:                preview.DiscountMessage,
		OtherPendingCryptoSubscription: preview.OtherPendingCryptoSubscription,
		OtherActiveSubscriptionId:      preview.OtherActiveSubscriptionId,
		ApplyPromoCredit:               preview.ApplyPromoCredit,
	}, nil
}
