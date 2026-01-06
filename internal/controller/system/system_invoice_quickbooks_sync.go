package system

import (
	"context"
	"unibee/internal/logic/analysis/quickbooks"

	"unibee/api/system/invoice"
)

func (c *ControllerInvoice) QuickbooksSync(ctx context.Context, req *invoice.QuickbooksSyncReq) (res *invoice.QuickbooksSyncRes, err error) {
	quickbooks.UploadPaidInvoice(ctx, req.InvoiceId)
	return &invoice.QuickbooksSyncRes{}, nil
}
