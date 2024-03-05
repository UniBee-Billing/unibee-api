package merchant

import (
	"context"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Edit(ctx context.Context, req *gateway.EditReq) (res *gateway.EditRes, err error) {
	gateway2.EditGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayId, req.GatewayKey, req.GatewaySecret)
	return &gateway.EditRes{}, nil
}
