package subscription

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee-api/internal/consumer/webhook/event"
	"unibee-api/internal/consumer/webhook/message"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

func printPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err)
	}
}

func SendSubscriptionCreatedMerchantWebhook(one *entity.Subscription) {
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
				printPanic(ctx, err)
				return
			}
		}()
		// todo mark need Construct SubscriptionDetail Data
		message.SendWebhookMessage(ctx, event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED, one.MerchantId, utility.FormatToGJson(one))
	}()
}
