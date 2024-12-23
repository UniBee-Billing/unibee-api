package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway/service"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) WireTransferSetup(ctx context.Context, req *gateway.WireTransferSetupReq) (res *gateway.WireTransferSetupRes, err error) {
	gateway2.SetupWireTransferGateway(ctx, &gateway2.WireTransferSetupReq{
		MerchantId:    _interface.GetMerchantId(ctx),
		Currency:      req.Currency,
		MinimumAmount: req.MinimumAmount,
		Bank:          req.Bank,
	})
	return &gateway.WireTransferSetupRes{}, nil
}
