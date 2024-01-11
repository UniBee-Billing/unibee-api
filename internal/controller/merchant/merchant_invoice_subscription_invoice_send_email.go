package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/handler"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceSendEmail(ctx context.Context, req *invoice.SubscriptionInvoiceSendEmailReq) (res *invoice.SubscriptionInvoiceSendEmailRes, err error) {
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SubscriptionInvoiceSendEmailRes{}, nil
}
