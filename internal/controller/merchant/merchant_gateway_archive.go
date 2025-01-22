package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) Archive(ctx context.Context, req *gateway.ArchiveReq) (res *gateway.ArchiveRes, err error) {
	return &gateway.ArchiveRes{Gateway: detail.ConvertGatewayDetail(ctx, service.ArchiveGateway(ctx, _interface.GetMerchantId(ctx), req.GatewayId))}, nil
}
