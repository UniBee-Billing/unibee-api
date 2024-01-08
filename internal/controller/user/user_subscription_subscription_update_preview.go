package user

import (
	"context"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error) {
	prepare, err := service.SubscriptionUpdatePreview(ctx, req, 0)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdatePreviewRes{
		TotalAmount: prepare.TotalAmount,
		Currency:    prepare.Currency,
		Invoice: &ro.ChannelDetailInvoiceRo{
			TotalAmount:        prepare.Invoice.TotalAmount,
			Currency:           prepare.Invoice.Currency,
			TaxAmount:          prepare.Invoice.TaxAmount,
			SubscriptionAmount: prepare.Invoice.SubscriptionAmount,
			Lines:              prepare.Invoice.Lines,
		},
		ProrationDate: prepare.ProrationDate,
	}, nil
}
