package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/query"

	"unibee/api/user/gateway"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx))
	return &gateway.ListRes{
		Gateways: ro.SimplifyGatewayList(data),
	}, nil
}
