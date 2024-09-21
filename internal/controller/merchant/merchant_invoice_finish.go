package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	"unibee/internal/cmd/i18n"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"
)

func (c *ControllerInvoice) Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Finish:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}
	return service.FinishInvoice(ctx, req)
}
