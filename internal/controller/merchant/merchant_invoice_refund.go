package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"
)

func (c *ControllerInvoice) Refund(ctx context.Context, req *invoice.RefundReq) (res *invoice.RefundRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Refund:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	refund, err := service.CreateInvoiceRefund(ctx, req)
	if err != nil {
		return nil, err
	}
	return &invoice.RefundRes{Refund: bean.SimplifyRefund(refund)}, nil
}
