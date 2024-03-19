package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	"unibee/internal/logic/invoice/handler"
	"unibee/utility"
)

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Send-Email:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	err = handler.SendSubscriptionInvoiceEmailToUser(ctx, req.InvoiceId)
	if err != nil {
		return nil, err
	}
	return &invoice.SendEmailRes{}, nil
}
