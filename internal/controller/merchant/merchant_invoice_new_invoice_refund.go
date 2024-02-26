package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) NewInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (res *invoice.NewInvoiceRefundRes, err error) {
	refund, err := service.CreateInvoiceRefund(ctx, req)
	if err != nil {
		return nil, err
	}
	return &invoice.NewInvoiceRefundRes{Refund: refund}, nil
}
