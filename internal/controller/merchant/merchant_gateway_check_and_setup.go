package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) CheckAndSetup(ctx context.Context, req *gateway.CheckAndSetupReq) (res *gateway.CheckAndSetupRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx)
	return &gateway.CheckAndSetupRes{}, nil
}
