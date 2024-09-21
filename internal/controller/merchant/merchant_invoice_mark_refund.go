package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/logic/invoice/service"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) MarkRefund(ctx context.Context, req *invoice.MarkRefundReq) (res *invoice.MarkRefundRes, err error) {
	refund, err := service.MarkInvoiceRefund(ctx, req)
	if err != nil {
		return nil, err
	}
	return &invoice.MarkRefundRes{Refund: bean.SimplifyRefund(refund)}, nil
}
