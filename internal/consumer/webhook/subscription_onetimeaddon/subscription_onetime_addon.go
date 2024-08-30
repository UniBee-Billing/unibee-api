package subscription_onetimeaddon

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean/detail"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	"unibee/utility"
)

func SendMerchantSubscriptionOnetimeAddonWebhookBackground(merchantId uint64, one *detail.SubscriptionOnetimeAddonDetail, event event.WebhookEvent) {
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
		message.SendWebhookMessage(ctx, event, merchantId, utility.FormatToGJson(one), "", "")
	}()
}
