package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean/detail"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/log"
	"unibee/internal/consumer/webhook/message"
	"unibee/internal/logic/metric_event"
	"unibee/internal/query"
	"unibee/utility"
)

func SendMerchantUserMetricWebhookBackground(userId uint64, sourceSubscriptionId string, event event.WebhookEvent, description string) {
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
		user := query.GetUserAccountById(ctx, userId)
		sub := query.GetSubscriptionBySubscriptionId(ctx, sourceSubscriptionId)
		if user != nil && sub != nil {
			plan := query.GetPlanById(ctx, sub.PlanId)
			if plan != nil {
				userMetric := metric_event.GetUserMetricStat(ctx, user.MerchantId, user, plan.ProductId, false)
				utility.AssertError(err, "SendMerchantUserMetricWebhookBackground Error")
				userMetric.Description = description

				message.SendWebhookMessage(ctx, event, user.MerchantId, utility.FormatToGJson(userMetric), "", "", nil)
			}
		}
	}()
}

func SendMerchantUserWebhookBackground(userId uint64, event event.WebhookEvent) {
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
		user := query.GetUserAccountById(ctx, userId)

		message.SendWebhookMessage(ctx, event, user.MerchantId, utility.FormatToGJson(detail.ConvertUserAccountToDetail(ctx, user)), "", "", nil)

	}()
}
