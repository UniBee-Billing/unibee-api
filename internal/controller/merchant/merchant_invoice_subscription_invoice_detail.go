package merchant

import (
	"context"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) SubscriptionInvoiceDetail(ctx context.Context, req *invoice.SubscriptionInvoiceDetailReq) (res *invoice.SubscriptionInvoiceDetailRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "InvoiceId Invalid")
	in := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(in != nil, "invoice not found")

	return &invoice.SubscriptionInvoiceDetailRes{Invoice: invoice_compute.ConvertInvoiceToRo(ctx, in)}, nil
}
