package merchant

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility/unibee"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	if !_interface.Context().Get(ctx).IsAdminPortalCall && req.Archive == nil {
		req.Archive = unibee.Bool(false)
	}
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx), req.Archive)

	gateways := gateway2.ConvertGatewayList(ctx, data)
	if _interface.Context().Get(ctx).IsOpenApiCall {
		for i, _ := range gateways {
			gateways[i].Bank = nil
		}
	}
	return &gateway.ListRes{
		Gateways: gateways,
	}, nil
}
