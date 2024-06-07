package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerInvoice) PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "invalid InvoiceId")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "invalid MerchantId")
	redisKey := fmt.Sprintf("Merchant-Invoice-PdfGenerate:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail, true)
	return &invoice.PdfGenerateRes{}, nil
}
