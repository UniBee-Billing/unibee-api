package user

import (
	"context"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (res *subscription.SubscriptionCreatePreviewRes, err error) {
	prepare, err := service.SubscriptionCreatePreview(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCreatePreviewRes{
		Plan:              prepare.Plan,
		Quantity:          prepare.Quantity,
		PayChannel:        service.ConvertChannelToRo(prepare.PayChannel),
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
