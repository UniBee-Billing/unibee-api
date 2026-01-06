package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) SetDefault(ctx context.Context, req *gateway.SetDefaultReq) (res *gateway.SetDefaultRes, err error) {
	return &gateway.SetDefaultRes{Gateway: detail.ConvertGatewayDetail(ctx, service.SetDefaultGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayId))}, nil
}
