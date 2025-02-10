package user

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	"unibee/api/user/gateway"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx), req.Archive)
	return &gateway.ListRes{
		Gateways: gateway2.ConvertGatewayList(ctx, data),
	}, nil
}
