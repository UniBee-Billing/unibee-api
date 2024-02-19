package merchant

import (
	"context"
	"unibee-api/internal/consumer/webhook/event"

	"unibee-api/api/merchant/webhook"
)

func (c *ControllerWebhook) EventList(ctx context.Context, req *webhook.EventListReq) (res *webhook.EventListRes, err error) {
	return &webhook.EventListRes{Events: event.ListeningEvents}, nil
}
