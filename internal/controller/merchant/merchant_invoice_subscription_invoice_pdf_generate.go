package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/handler"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoicePdfGenerate(ctx context.Context, req *invoice.SubscriptionInvoicePdfGenerateReq) (res *invoice.SubscriptionInvoicePdfGenerateRes, err error) {
	_ = handler.SubscriptionInvoicePdfGenerateBackground(req.InvoiceId)
	return &invoice.SubscriptionInvoicePdfGenerateRes{}, nil
}
