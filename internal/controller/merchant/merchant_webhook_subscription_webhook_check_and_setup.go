package merchant

import (
	"context"
	"go-oversea-pay/api/merchant/webhook"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/channel"
	"go-oversea-pay/utility"
)

func (c *ControllerWebhook) SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	channel.CheckAndSetupPayChannelWebhooks(ctx)
	return &webhook.SubscriptionWebhookCheckAndSetupRes{}, nil
}
