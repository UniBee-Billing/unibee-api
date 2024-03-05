package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) New(ctx context.Context, req *invoice.NewReq) (res *invoice.NewRes, err error) {
	return service.CreateInvoice(ctx, req)
}
