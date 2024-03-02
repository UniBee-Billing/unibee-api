package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) Cancel(ctx context.Context, req *invoice.CancelReq) (res *invoice.CancelRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	return &invoice.CancelRes{}, service.CancelProcessingInvoice(ctx, req.InvoiceId)
}
