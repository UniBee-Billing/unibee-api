package payment

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	"unibee/internal/logic/payment/detail"
	"unibee/internal/query"
	"unibee/utility"
)

func SendPaymentWebhookBackground(paymentId string, event event.WebhookEvent) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		one := query.GetPaymentByPaymentId(ctx, paymentId)
		if one != nil {
			paymentDetail := detail.GetPaymentDetail(ctx, one.MerchantId, one.PaymentId)
			utility.Assert(paymentDetail != nil, "SendPaymentWebhookBackground Error")

			message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(paymentDetail), "", "")
		}
	}()
}
