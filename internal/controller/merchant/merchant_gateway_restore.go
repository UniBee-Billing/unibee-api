package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/api/merchant/gateway"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/service"
)

func (c *ControllerGateway) Restore(ctx context.Context, req *gateway.RestoreReq) (res *gateway.RestoreRes, err error) {
	return &gateway.RestoreRes{Gateway: detail.ConvertGatewayDetail(ctx, service.RestoreGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayId))}, nil
}
