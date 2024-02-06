package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/invoice/service"
	"unibee-api/utility"

	"unibee-api/api/merchant/invoice"
)

func (c *ControllerInvoice) DeletePendingInvoice(ctx context.Context, req *invoice.DeletePendingInvoiceReq) (res *invoice.DeletePendingInvoiceRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "merchantUserId invalid")
		//utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId == uint64(req.MerchantId), "merchantId not match")
	}

	return &invoice.DeletePendingInvoiceRes{}, service.DeletePendingInvoice(ctx, req.InvoiceId)
}
