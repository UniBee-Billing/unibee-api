package user

import (
	"context"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/subscription/service"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	prepare, err := service.SubscriptionCreatePreview(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.CreatePreviewRes{
		Plan:              prepare.Plan,
		Quantity:          prepare.Quantity,
		Gateway:           service.ConvertGatewayToRo(prepare.Gateway),
		AddonParams:       prepare.AddonParams,
		Addons:            prepare.Addons,
		TotalAmount:       prepare.TotalAmount,
		Currency:          prepare.Currency,
		VatNumber:         prepare.VatNumber,
		VatNumberValidate: prepare.VatNumberValidate,
		VatCountryCode:    prepare.VatCountryCode,
		VatCountryName:    prepare.VatCountryName,
		TaxScale:          prepare.TaxScale,
		Invoice: &ro.InvoiceDetailSimplify{
			InvoiceName:                    prepare.Invoice.InvoiceName,
			TotalAmount:                    prepare.Invoice.TotalAmount,
			TotalAmountExcludingTax:        prepare.Invoice.TotalAmountExcludingTax,
			Currency:                       prepare.Invoice.Currency,
			TaxAmount:                      prepare.Invoice.TaxAmount,
			SubscriptionAmount:             prepare.Invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: prepare.Invoice.SubscriptionAmountExcludingTax,
			Lines:                          prepare.Invoice.Lines,
		},
		UserId: prepare.UserId,
		Email:  prepare.Email,
	}, nil
}
