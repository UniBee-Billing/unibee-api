package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/handler"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	err = handler.SendSubscriptionInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SendEmailRes{}, nil
}