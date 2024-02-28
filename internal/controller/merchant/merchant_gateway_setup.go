package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway/service"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	gateway2.SetupGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.GatewayKey, req.GatewaySecret)
	return &gateway.SetupRes{}, nil
}
