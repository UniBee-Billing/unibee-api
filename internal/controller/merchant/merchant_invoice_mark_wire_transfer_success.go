package merchant

import (
	"context"
	"unibee/internal/logic/invoice/handler"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) MarkWireTransferSuccess(ctx context.Context, req *invoice.MarkWireTransferSuccessReq) (res *invoice.MarkWireTransferSuccessRes, err error) {
	_, err = handler.MarkWireTransferInvoiceAsSuccess(ctx, req.InvoiceId, req.TransferNumber)
	if err != nil {
		return nil, err
	}
	return &invoice.MarkWireTransferSuccessRes{}, nil
}
