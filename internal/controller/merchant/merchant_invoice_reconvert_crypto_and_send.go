package merchant

import (
	"context"
	"unibee/internal/logic/invoice/handler"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) ReconvertCryptoAndSend(ctx context.Context, req *invoice.ReconvertCryptoAndSendReq) (res *invoice.ReconvertCryptoAndSendRes, err error) {
	err = handler.ReconvertCryptoDataForInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	err = handler.SendSubscriptionInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
