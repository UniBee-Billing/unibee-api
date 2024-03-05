package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) Edit(ctx context.Context, req *invoice.EditReq) (res *invoice.EditRes, err error) {
	return service.EditInvoice(ctx, req)
}
