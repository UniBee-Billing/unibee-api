package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferEdit(ctx context.Context, req *gateway.WireTransferEditReq) (res *gateway.WireTransferEditRes, err error) {
	one := gateway2.EditWireTransferGateway(ctx, &gateway2.WireTransferSetupReq{
		GatewayId:     req.GatewayId,
		MerchantId:    _interface.GetMerchantId(ctx),
		Currency:      req.Currency,
		MinimumAmount: req.MinimumAmount,
		Bank:          req.Bank,
		DisplayName:   req.DisplayName,
		GatewayIcon:   req.GatewayIcons,
		Sort:          req.Sort,
	})
	return &gateway.WireTransferEditRes{Gateway: detail.ConvertGatewayDetail(ctx, one)}, nil
}
