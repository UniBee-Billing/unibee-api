package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/invoice_compute"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "InvoiceId Invalid")
	in := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(in != nil, "invoice not found")
	utility.Assert(in.MerchantId == _interface.GetMerchantId(ctx), "wrong merchant account")

	return &invoice.DetailRes{Invoice: invoice_compute.ConvertInvoiceToRo(ctx, in)}, nil
}
