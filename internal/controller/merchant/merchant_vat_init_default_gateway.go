package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/vat_gateway"

	"unibee/api/merchant/vat"
)

func (c *ControllerVat) InitDefaultGateway(ctx context.Context, req *vat.InitDefaultGatewayReq) (res *vat.InitDefaultGatewayRes, err error) {
	err = vat_gateway.InitMerchantDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
	if err != nil {
		return nil, err
	}
	return &vat.InitDefaultGatewayRes{}, nil
}
