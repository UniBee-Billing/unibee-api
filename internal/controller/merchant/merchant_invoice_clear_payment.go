package merchant

import (
	"context"
	"fmt"
	"unibee/internal/consts"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/invoice"
)

func (c *ControllerInvoice) ClearPayment(ctx context.Context, req *invoice.ClearPaymentReq) (res *invoice.ClearPaymentRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "Invalid InvoiceId")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "Invoice not found")
	utility.Assert(one.Status == consts.InvoiceStatusProcessing, "Invoice not processing")
	_, err = service.ClearInvoicePayment(ctx, one)
	if err != nil {
		return nil, err
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "ClearInvoicePaymentBYAdmin",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &invoice.ClearPaymentRes{}, nil
}
