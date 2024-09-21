package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) Cancel(ctx context.Context, req *invoice.CancelReq) (res *invoice.CancelRes, err error) {
	return &invoice.CancelRes{}, service.CancelProcessingInvoice(ctx, req.InvoiceId, "AdminCancel")
}
