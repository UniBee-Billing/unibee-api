package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) Delete(ctx context.Context, req *invoice.DeleteReq) (res *invoice.DeleteRes, err error) {
	return &invoice.DeleteRes{}, service.DeletePendingInvoice(ctx, req.InvoiceId)
}
