package merchant

import (
	"context"
	"unibee/api/merchant/vat"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/vat_gateway/setup"
	"unibee/utility"
)

func (c *ControllerVat) SetupGateway(ctx context.Context, req *vat.SetupGatewayReq) (res *vat.SetupGatewayRes, err error) {
	err = setup.SetupMerchantVatConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.Data, req.IsDefault)
	utility.AssertError(err, "Setup Vat Gateway Error")
	if req.IsDefault {
		err = setup.InitMerchantDefaultVatGateway(ctx, _interface.GetMerchantId(ctx))
		utility.AssertError(err, "Init Vat Gateway Error")
	}
	return &vat.SetupGatewayRes{Data: utility.HideStar(req.Data)}, nil
}
