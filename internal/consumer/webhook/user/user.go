package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	"unibee/internal/logic/metric_event"
	"unibee/internal/query"
	"unibee/utility"
)

func SendMerchantUserMetricWebhookBackground(userId int64, event event.MerchantWebhookEvent) {
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
		user := query.GetUserAccountById(ctx, uint64(userId))
		if user != nil {
			userMetric := metric_event.GetUserMetricStat(ctx, user.MerchantId, user)
			utility.AssertError(err, "SendMerchantUserMetricWebhookBackground Error")

			message.SendWebhookMessage(ctx, event, user.MerchantId, utility.FormatToGJson(userMetric))
		}
	}()
}
