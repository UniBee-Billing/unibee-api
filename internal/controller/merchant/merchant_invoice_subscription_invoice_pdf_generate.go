package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/invoice"
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
