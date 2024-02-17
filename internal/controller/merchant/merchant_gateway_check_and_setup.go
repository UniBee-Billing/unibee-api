package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	gatewayWebhook "unibee-api/internal/logic/gateway/webhook"
	"unibee-api/utility"

	"unibee-api/api/merchant/gateway"
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
