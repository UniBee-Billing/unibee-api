package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/handler"
)

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {
	err = handler.SendSubscriptionInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SendEmailRes{}, nil
}
