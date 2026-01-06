package system

import (
	"context"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	subscription3 "unibee/internal/consumer/webhook/subscription"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/subscription"
)

func (c *ControllerSubscription) BatchSendSubActivateWebhookEvent(ctx context.Context, req *subscription.BatchSendSubActivateWebhookEventReq) (res *subscription.BatchSendSubActivateWebhookEventRes, err error) {
	utility.Assert(len(req.SubIds) > 0, "Empty SubIds")
	for _, subId := range req.SubIds {
		sub := query.GetSubscriptionBySubscriptionId(ctx, subId)
		if sub != nil {
			if sub.Status == consts.SubStatusActive {
				subscription3.SendMerchantSubscriptionWebhookBackground(sub, -10000, event.UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_ACTIVATED, map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()})
			}
		}
	}
	return &subscription.BatchSendSubActivateWebhookEventRes{}, nil
}
