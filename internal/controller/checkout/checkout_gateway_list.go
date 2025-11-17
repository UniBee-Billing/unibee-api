package checkout

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	"unibee/api/checkout/gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	data := query.GetMerchantGatewayList(ctx, req.MerchantId, nil)
	list := make([]*entity.MerchantGateway, 0)
	for _, item := range data {
		//if item.GatewayType != consts.GatewayTypeWireTransfer {
		list = append(list, item)
		//}
	}
	return &gateway.ListRes{
		Gateways: gateway2.ConvertGatewayList(ctx, list),
	}, nil
}
