package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"
)

func (c *ControllerInvoice) Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Finish:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	return service.FinishInvoice(ctx, req)
}
