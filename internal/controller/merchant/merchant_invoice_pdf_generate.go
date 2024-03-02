package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/handler"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail)
	return &invoice.PdfGenerateRes{}, nil
}
