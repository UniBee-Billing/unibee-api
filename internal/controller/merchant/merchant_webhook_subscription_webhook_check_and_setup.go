package merchant

import (
	"context"
	"unibee-api/api/merchant/webhook"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	webhook2 "unibee-api/internal/logic/gateway/webhook"
	"unibee-api/utility"
)

func (c *ControllerWebhook) SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	webhook2.CheckAndSetupGatewayWebhooks(ctx)
	return &webhook.SubscriptionWebhookCheckAndSetupRes{}, nil
}
