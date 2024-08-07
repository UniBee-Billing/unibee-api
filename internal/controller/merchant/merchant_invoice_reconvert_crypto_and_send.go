package merchant

import (
	"context"
	"fmt"
	"unibee/internal/cmd/i18n"
	"unibee/internal/logic/invoice/handler"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) ReconvertCryptoAndSend(ctx context.Context, req *invoice.ReconvertCryptoAndSendReq) (res *invoice.ReconvertCryptoAndSendRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-ReconvertCryptoAndSend:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}
	err = handler.ReconvertCryptoDataForInvoice(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId, true)
	if err != nil {
		return nil, err
	}
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
