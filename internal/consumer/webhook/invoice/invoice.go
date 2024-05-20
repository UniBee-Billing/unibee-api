package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean/detail"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func SendMerchantInvoiceWebhookBackground(one *entity.Invoice, event event.MerchantWebhookEvent) {
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
		key := fmt.Sprintf("webhook_lock_%s_%s", one.InvoiceId, event)
		if !utility.TryLock(ctx, key, 10*60) {
			message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(detail.ConvertInvoiceToDetail(ctx, one)))
		}
	}()
}
