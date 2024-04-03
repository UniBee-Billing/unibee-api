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

func SendRefundWebhookBackground(refundId string, event event.MerchantWebhookEvent) {
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
		one := query.GetRefundByRefundId(ctx, refundId)
		if one != nil {
			refundDetail := detail.GetRefundDetail(ctx, one.MerchantId, one.RefundId)
			utility.Assert(refundDetail != nil, "SendRefundWebhookBackground Error")

			message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(refundDetail))
		}
	}()
}
