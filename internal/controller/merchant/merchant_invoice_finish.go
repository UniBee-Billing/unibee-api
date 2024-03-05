package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error) {
	return service.FinishInvoice(ctx, req)
}
