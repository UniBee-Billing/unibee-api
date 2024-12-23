package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error) {
	return &gateway.SetupRes{Gateway: bean.SimplifyGateway(gateway2.SetupGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.GatewayKey, req.GatewaySecret))}, nil
}
