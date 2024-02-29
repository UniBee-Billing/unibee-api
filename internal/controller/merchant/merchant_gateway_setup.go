package merchant

import (
	"context"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error) {
	gateway2.SetupGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.GatewayKey, req.GatewaySecret)
	return &gateway.SetupRes{}, nil
}
