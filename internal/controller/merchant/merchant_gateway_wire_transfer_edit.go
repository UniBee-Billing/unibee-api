package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	gateway2 "unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferEdit(ctx context.Context, req *gateway.WireTransferEditReq) (res *gateway.WireTransferEditRes, err error) {
	gateway2.EditWireTransferGateway(ctx, &gateway2.WireTransferSetupReq{
		GatewayId:     req.GatewayId,
		MerchantId:    _interface.GetMerchantId(ctx),
		Currency:      req.Currency,
		MinimumAmount: req.MinimumAmount,
		Bank:          req.Bank,
	})
	return &gateway.WireTransferEditRes{}, nil
}
