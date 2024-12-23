package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) New(ctx context.Context, req *invoice.NewReq) (res *invoice.NewRes, err error) {
	return service.CreateInvoice(ctx, _interface.GetMerchantId(ctx), req)
}
