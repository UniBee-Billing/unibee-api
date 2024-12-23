package merchant

import (
	"context"
	"unibee/api/merchant/email"
	_interface "unibee/internal/interface/context"
	email2 "unibee/internal/logic/email"
)

func (c *ControllerEmail) GatewaySetup(ctx context.Context, req *email.GatewaySetupReq) (res *email.GatewaySetupRes, err error) {
	err = email2.SetupMerchantEmailConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.Data, req.IsDefault)
	if err != nil {
		return nil, err
	}
	return &email.GatewaySetupRes{}, nil
}
