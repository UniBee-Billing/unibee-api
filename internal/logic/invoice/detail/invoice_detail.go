package detail

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/internal/query"
)

func InvoiceDetail(ctx context.Context, invoiceId string) *detail.InvoiceDetail {
	one := query.GetInvoiceByInvoiceId(ctx, invoiceId)
	if one != nil {
		return detail.ConvertInvoiceToDetail(ctx, one)
	}
	return nil
}
