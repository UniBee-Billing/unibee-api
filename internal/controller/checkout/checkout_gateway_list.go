package checkout

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	"unibee/api/checkout/gateway"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	data := query.GetMerchantGatewayList(ctx, req.MerchantId, nil)
	return &gateway.ListRes{
		Gateways: gateway2.ConvertGatewayList(ctx, data),
	}, nil
}
