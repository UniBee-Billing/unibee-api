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

	return &gateway.ListRes{
		Gateways: gateway2.ConvertGatewayList(ctx, data),
	}, nil
}
