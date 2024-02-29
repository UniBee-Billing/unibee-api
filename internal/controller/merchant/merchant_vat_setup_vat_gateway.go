package merchant

import (
	"context"
	"unibee/api/merchant/vat"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/vat_gateway"
)

func (c *ControllerVat) SetupVatGateway(ctx context.Context, req *vat.SetupVatGatewayReq) (res *vat.SetupVatGatewayRes, err error) {
	err = vat_gateway.SetupMerchantVatConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.Data, req.IsDefault)
	if err != nil {
		return nil, err
	}
	if req.IsDefault {
		err := vat_gateway.InitMerchantDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
		if err != nil {
			return nil, err
		}
	}
	return &vat.SetupVatGatewayRes{}, nil
}
