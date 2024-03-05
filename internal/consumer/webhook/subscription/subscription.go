package subscription

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/message"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func printPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err)
	}
}

func SendMerchantSubscriptionWebhookBackground(one *entity.Subscription, event event.MerchantWebhookEvent) {
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
		subDetailRes, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
		utility.AssertError(err, "SendMerchantSubscriptionWebhookBackground SubscriptionDetail Error")

		message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(subDetailRes))
	}()
}
