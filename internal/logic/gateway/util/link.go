package util

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"strings"
	entity "unibee/internal/model/entity/default"
)

func GetPaymentRedirectUrl(ctx context.Context, payment *entity.Payment, success string) string {
	var targetUrl = ""
	if payment == nil {
		return targetUrl
	}
	if success == "false" {
		var metadata = make(map[string]string)
		if len(payment.MetaData) > 0 {
			err := gjson.Unmarshal([]byte(payment.MetaData), &metadata)
			if err != nil {
				fmt.Printf("SimplifyPayment Unmarshal Metadata error:%s", err.Error())
			}
		}
		cancelUrl := metadata["CancelUrl"]
		if cancelUrl != "" && len(cancelUrl) > 0 {
			targetUrl = cancelUrl
		} else {
			targetUrl = payment.ReturnUrl
		}
	} else {
		targetUrl = payment.ReturnUrl
	}
	if len(targetUrl) == 0 {
		return targetUrl
	}
	if strings.Contains(targetUrl, "?") {
		targetUrl = targetUrl + fmt.Sprintf("&paymentId=%s&subId=%s&invoiceId=%s&success=%v", payment.PaymentId, payment.SubscriptionId, payment.InvoiceId, success)
	} else {
		targetUrl = targetUrl + fmt.Sprintf("?paymentId=%s&subId=%s&invoiceId=%s&success=%v", payment.PaymentId, payment.SubscriptionId, payment.InvoiceId, success)
	}
	return targetUrl
}
