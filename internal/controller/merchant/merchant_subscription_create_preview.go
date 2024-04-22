package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	prepare, err := service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		PlanId:         req.PlanId,
		UserId:         req.UserId,
		Quantity:       req.Quantity,
		GatewayId:      req.GatewayId,
		AddonParams:    req.AddonParams,
		VatCountryCode: req.VatCountryCode,
		VatNumber:      req.VatNumber,
		TaxPercentage:  req.TaxPercentage,
		DiscountCode:   req.DiscountCode,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.CreatePreviewRes{
		Plan:              bean.SimplifyPlan(prepare.Plan),
		TrialEnd:          prepare.TrialEnd,
		Quantity:          prepare.Quantity,
		Gateway:           bean.SimplifyGateway(prepare.Gateway),
		AddonParams:       prepare.AddonParams,
		Addons:            prepare.Addons,
		OriginAmount:      prepare.OriginAmount,
		TotalAmount:       prepare.TotalAmount,
		DiscountAmount:    prepare.DiscountAmount,
		Currency:          prepare.Currency,
		VatNumber:         prepare.VatNumber,
		VatNumberValidate: prepare.VatNumberValidate,
		VatCountryCode:    prepare.VatCountryCode,
		VatCountryName:    prepare.VatCountryName,
		TaxPercentage:     prepare.TaxPercentage,
		Invoice: &bean.InvoiceSimplify{
			InvoiceName:                    prepare.Invoice.InvoiceName,
			TotalAmount:                    prepare.Invoice.TotalAmount,
			TotalAmountExcludingTax:        prepare.Invoice.TotalAmountExcludingTax,
			Currency:                       prepare.Invoice.Currency,
			TaxAmount:                      prepare.Invoice.TaxAmount,
			SubscriptionAmount:             prepare.Invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: prepare.Invoice.SubscriptionAmountExcludingTax,
			Lines:                          prepare.Invoice.Lines,
		},
		UserId:   prepare.UserId,
		Email:    prepare.Email,
		Discount: prepare.Discount,
	}, nil
}
