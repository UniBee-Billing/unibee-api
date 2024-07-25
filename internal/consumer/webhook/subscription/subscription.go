package subscription

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func SendMerchantSubscriptionWebhookBackground(one *entity.Subscription, dayLeft int, event event.MerchantWebhookEvent) {
	go func() {
		ctx := context.Background()
		g.Log().Infof(ctx, "SendMerchantSubscriptionWebhookBackground event:%v", event)
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
		subDetailRes, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
		utility.AssertError(err, "SendMerchantSubscriptionWebhookBackground SubscriptionDetail Error")
		if dayLeft == -10000 {
			subDetailRes.DayLeft = int((subDetailRes.Subscription.CurrentPeriodEnd - utility.MaxInt64(gtime.Now().Timestamp(), subDetailRes.Subscription.TestClock) + 7200) / 86400)
		} else {
			subDetailRes.DayLeft = dayLeft
		}
		message.SendWebhookMessage(ctx, event, one.MerchantId, utility.FormatToGJson(subDetailRes), "", "")
	}()
}
