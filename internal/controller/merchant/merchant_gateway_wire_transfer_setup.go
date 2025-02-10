package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferSetup(ctx context.Context, req *gateway.WireTransferSetupReq) (res *gateway.WireTransferSetupRes, err error) {
	one := gateway2.SetupWireTransferGateway(ctx, &gateway2.WireTransferSetupReq{
		MerchantId:    _interface.GetMerchantId(ctx),
		Currency:      req.Currency,
		MinimumAmount: req.MinimumAmount,
		Bank:          req.Bank,
		DisplayName:   req.DisplayName,
		GatewayIcon:   req.GatewayIcons,
		Sort:          req.Sort,
	})
	return &gateway.WireTransferSetupRes{Gateway: detail.ConvertGatewayDetail(ctx, one)}, nil
}
