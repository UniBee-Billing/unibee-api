package merchant

import (
	"context"
	"unibee/internal/consumer/webhook/message"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) ResendWebhook(ctx context.Context, req *webhook.ResendWebhookReq) (res *webhook.ResendWebhookRes, err error) {
	return &webhook.ResendWebhookRes{SendResult: message.ResentWebhook(ctx, req.LogId)}, nil
}
