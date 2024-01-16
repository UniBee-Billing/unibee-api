package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/handler"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceSendEmail(ctx context.Context, req *invoice.SubscriptionInvoiceSendEmailReq) (res *invoice.SubscriptionInvoiceSendEmailRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SubscriptionInvoiceSendEmailRes{}, nil
}
