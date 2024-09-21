package invoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func SendMerchantInvoiceWebhookBackground(one *entity.Invoice, event event.WebhookEvent, metadata map[string]interface{}) {
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
		g.Log().Infof(ctx, "SendMerchantInvoiceWebhookBackground_invoiceId:%s， event:%s", one.InvoiceId, event)
		if one != nil {
			key := fmt.Sprintf("webhook_invoice_lock_%s_%s", one.InvoiceId, event)
			if utility.TryLock(ctx, key, 60) {
				message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(detail.ConvertInvoiceToDetail(ctx, one)), "", "", metadata)
			}
		}
	}()
}
