package checkout

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean"
	"unibee/internal/logic/gateway/util"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/checkout/payment"
)

func (c *ControllerPayment) Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	var targetUrl = util.GetPaymentRedirectUrl(ctx, one, "true")
	if len(targetUrl) == 0 {
		merchant := query.GetMerchantById(ctx, one.MerchantId)
		if merchant != nil && len(merchant.Host) > 0 {
			if strings.HasPrefix(merchant.Host, "http") {
				targetUrl = merchant.Host
			} else {
				targetUrl = fmt.Sprintf("http://%s", merchant.Host)
			}
		}
	}
	cancelUrl := util.GetPaymentRedirectUrl(ctx, one, "false")
	if len(cancelUrl) == 0 {
		merchant := query.GetMerchantById(ctx, one.MerchantId)
		if merchant != nil && len(merchant.Host) > 0 {
			if strings.HasPrefix(merchant.Host, "http") {
				cancelUrl = merchant.Host
			} else {
				cancelUrl = fmt.Sprintf("http://%s", merchant.Host)
			}
		}
	}

	return &payment.DetailRes{
		PaymentStatus: one.Status,
		Payment:       bean.SimplifyPayment(one),
		ReturnUrl:     targetUrl,
		CancelUrl:     cancelUrl,
	}, nil
}
