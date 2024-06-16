package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/invoice"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/member"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerInvoice) SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error) {
	redisKey := fmt.Sprintf("Merchant-Invoice-Send-Email:%s", req.InvoiceId)
	if !utility.TryLock(ctx, redisKey, 10) {
		utility.Assert(false, "click too fast, please wait for second")
	}
	one := query.GetInvoiceByInvoiceId(ctx, req.InvoiceId)
	utility.Assert(one != nil, "invoice not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "invalid MerchantId")
	err = handler.SendInvoiceEmailToUser(ctx, req.InvoiceId, true)
	member.AppendOptLog(ctx, &member.OptLogRequest{
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
