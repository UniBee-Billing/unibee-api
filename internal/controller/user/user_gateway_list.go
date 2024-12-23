package user

import (
	"context"
	gateway2 "unibee/api/bean"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"

	"unibee/api/user/gateway"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx))
	return &gateway.ListRes{
		Gateways: gateway2.SimplifyGatewayList(data),
	}, nil
}
