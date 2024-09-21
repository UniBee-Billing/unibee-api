package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/service"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) MarkRefundInvoiceSuccess(ctx context.Context, req *invoice.MarkRefundInvoiceSuccessReq) (res *invoice.MarkRefundInvoiceSuccessRes, err error) {
	service.MarkInvoiceRefundSuccess(ctx, _interface.GetMerchantId(ctx), req.InvoiceId, req.Reason)
	return &invoice.MarkRefundInvoiceSuccessRes{}, nil
}
