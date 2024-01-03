package webhook

import (
	"context"
	"go-oversea-pay/internal/logic/payment/outchannel"

	"go-oversea-pay/api/webhook/v1"
)

func (c *ControllerV1) SubscriptionWebhookCheckAndSetup(ctx context.Context, req *v1.SubscriptionWebhookCheckAndSetupReq) (res *v1.SubscriptionWebhookCheckAndSetupRes, err error) {
	outchannel.CheckAndSetupPayChannelWebhooks(ctx)
	return &v1.SubscriptionWebhookCheckAndSetupRes{}, nil
}
