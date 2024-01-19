package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/merchant/invoice"
)

func (c *ControllerInvoice) NewInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (res *invoice.NewInvoiceRefundRes, err error) {
	refund, err := service.CreateInvoiceRefund(ctx, req)
	if err != nil {
		return nil, err
	}
	return &invoice.NewInvoiceRefundRes{Refund: refund}, nil
}
