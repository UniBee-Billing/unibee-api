package merchant

import (
	"context"
	"unibee/api/bean/detail"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "InvoiceId Invalid")
	in := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(in != nil, "invoice not found")
	utility.Assert(in.MerchantId == _interface.GetMerchantId(ctx), "wrong merchant account")

	var creditNoteEntities []*entity.Invoice
	var creditNotes = make([]*detail.InvoiceDetail, 0)
	if len(in.RefundId) == 0 {
		_ = dao.Invoice.Ctx(ctx).
			Where(dao.Invoice.Columns().MerchantId, in.MerchantId).
			Where(dao.Invoice.Columns().PaymentId, in.PaymentId).
			WhereNotNull(dao.Invoice.Columns().RefundId).Scan(&creditNoteEntities)
		for _, one := range creditNoteEntities {
			if one.Id != in.Id {
				creditNotes = append(creditNotes, detail.ConvertInvoiceToDetail(ctx, one))
			}
		}
	}

	return &invoice.DetailRes{Invoice: detail.ConvertInvoiceToDetail(ctx, in), CreditNotes: creditNotes}, nil
}
