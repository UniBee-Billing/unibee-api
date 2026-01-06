package checkout

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	"unibee/api/checkout/gateway"
	"unibee/internal/consts"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	data := query.GetMerchantGatewayList(ctx, req.MerchantId, unibee.Bool(false))
	list := make([]*entity.MerchantGateway, 0)
	for _, item := range data {
		if item.GatewayType != consts.GatewayTypeWireTransfer {
			list = append(list, item)
		}
	}
	return &gateway.ListRes{
		Gateways: gateway2.ConvertGatewayList(ctx, list),
	}, nil
}
