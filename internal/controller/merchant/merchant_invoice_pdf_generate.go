package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	"unibee/internal/cmd/i18n"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/operation_log"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerInvoice) PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error) {
	utility.Assert(len(req.InvoiceId) > 0, "invalid InvoiceId")
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "invalid MerchantId")
	redisKey := fmt.Sprintf("Merchant-Invoice-PdfGenerate:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}
	_ = handler.InvoicePdfGenerateAndEmailSendBackground(req.InvoiceId, req.SendUserEmail, true)
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "PdfGenerateAndSend",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return &invoice.PdfGenerateRes{}, nil
}
