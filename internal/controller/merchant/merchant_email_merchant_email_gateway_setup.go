package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) MerchantEmailGatewaySetup(ctx context.Context, req *email.MerchantEmailGatewaySetupReq) (res *email.MerchantEmailGatewaySetupRes, err error) {
	err = email2.SetupMerchantEmailConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.Data, req.IsDefault)
	if err != nil {
		return nil, err
	}
	return &email.MerchantEmailGatewaySetupRes{}, nil
}
