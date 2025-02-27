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

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Send-Email:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, i18n.LocalizationFormat(ctx, "{#ClickTooFast}"))
	}
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "invalid MerchantId")
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId, true, "")
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Invoice(%s)", one.InvoiceId),
		Content:        "Send",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return nil, err
	}
	return &invoice.SendEmailRes{}, nil
}
