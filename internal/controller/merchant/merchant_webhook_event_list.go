package merchant

import (
	"context"
	"unibee/internal/consumer/webhook/event"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) EventList(ctx context.Context, req *webhook.EventListReq) (res *webhook.EventListRes, err error) {
	return &webhook.EventListRes{EventList: event.ListeningEventList}, nil
}
