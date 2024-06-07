package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/handler"
	"unibee/utility"
)

func (c *ControllerInvoice) PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-PdfGenerate:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail, true)
	return &invoice.PdfGenerateRes{}, nil
}
