package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/handler"
)

func (c *ControllerInvoice) PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error) {
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail)
	return &invoice.PdfGenerateRes{}, nil
}
