package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) Detail(ctx context.Context, req *gateway.DetailReq) (res *gateway.DetailRes, err error) {
	utility.Assert(req.GatewayId != nil || req.GatewayName != nil, "Either gatewayId or gatewayName needed")
	if req.GatewayId != nil {
		one := query.GetGatewayById(ctx, *req.GatewayId)
		utility.Assert(one != nil, "Gateway not found")
		return &gateway.DetailRes{Gateway: detail.ConvertGatewayDetail(ctx, one)}, nil
	} else {
		one := query.GetGatewayByGatewayName(ctx, _interface.GetMerchantId(ctx), *req.GatewayName)
		utility.Assert(one != nil, "Gateway not setup")
		return &gateway.DetailRes{Gateway: detail.ConvertGatewayDetail(ctx, one)}, nil
	}
}
