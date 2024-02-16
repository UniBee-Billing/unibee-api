package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/invoice/handler"
	"unibee-api/utility"

	"unibee-api/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceSendEmail(ctx context.Context, req *invoice.SubscriptionInvoiceSendEmailReq) (res *invoice.SubscriptionInvoiceSendEmailRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	err = handler.SendSubscriptionInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SubscriptionInvoiceSendEmailRes{}, nil
}
