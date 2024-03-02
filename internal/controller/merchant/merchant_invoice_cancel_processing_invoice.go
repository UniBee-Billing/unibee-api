package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) CancelProcessingInvoice(ctx context.Context, req *invoice.CancelProcessingInvoiceReq) (res *invoice.CancelProcessingInvoiceRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	return &invoice.CancelProcessingInvoiceRes{}, service.CancelProcessingInvoice(ctx, req.InvoiceId)
}
