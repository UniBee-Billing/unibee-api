package merchant

import (
	"context"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
)

func (c *ControllerInvoice) Refund(ctx context.Context, req *invoice.RefundReq) (res *invoice.RefundRes, err error) {
	refund, err := service.CreateInvoiceRefund(ctx, req)
	if err != nil {
		return nil, err
	}
	return &invoice.RefundRes{Refund: refund}, nil
}
