package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/invoice/handler"
	"unibee-api/utility"

	"unibee-api/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoicePdfGenerate(ctx context.Context, req *invoice.SubscriptionInvoicePdfGenerateReq) (res *invoice.SubscriptionInvoicePdfGenerateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	_ = handler.SubscriptionInvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail)
	return &invoice.SubscriptionInvoicePdfGenerateRes{}, nil
}
