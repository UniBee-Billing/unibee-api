package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/invoice/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) NewInvoiceCreate(ctx context.Context, req *invoice.NewInvoiceCreateReq) (res *invoice.NewInvoiceCreateRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId == uint64(req.MerchantId), "merchantId not match")
	}

	return service.CreateInvoice(ctx, req)
}
