package user

import (
	"context"
	gateway2 "unibee/api/bean/detail"
	"unibee/api/user/gateway"
	_interface "unibee/internal/interface/context"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func (c *ControllerGateway) List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error) {
	data := query.GetMerchantGatewayList(ctx, _interface.GetMerchantId(ctx), req.Archive)
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
