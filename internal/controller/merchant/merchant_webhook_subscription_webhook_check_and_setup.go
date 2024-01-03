package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/webhook"
	"go-oversea-pay/internal/logic/payment/outchannel"
)

func (c *ControllerWebhook) SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error) {
	outchannel.CheckAndSetupPayChannelWebhooks(ctx)
	return &webhook.SubscriptionWebhookCheckAndSetupRes{}, nil
}
