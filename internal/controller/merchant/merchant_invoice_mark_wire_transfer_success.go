package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) MarkWireTransferSuccess(ctx context.Context, req *invoice.MarkWireTransferSuccessReq) (res *invoice.MarkWireTransferSuccessRes, err error) {
	_, err = service.MarkWireTransferInvoiceAsSuccess(ctx, req.InvoiceId, req.TransferNumber, req.Reason)
	if err != nil {
		return nil, err
	}
	return &invoice.MarkWireTransferSuccessRes{}, nil
}
