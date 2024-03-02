package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	return service.FinishInvoice(ctx, req)
}
