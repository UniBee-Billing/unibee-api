package merchant

import (
	"context"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "InvoiceId Invalid")
	in := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(in != nil, "invoice not found")

	return &invoice.DetailRes{Invoice: invoice_compute.ConvertInvoiceToRo(ctx, in)}, nil
}
