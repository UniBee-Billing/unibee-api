package merchant

import (
	"context"
	"unibee/api/merchant/email"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	email2 "unibee/internal/logic/email"
	"unibee/utility"
)

func (c *ControllerEmail) MerchantEmailGatewaySetup(ctx context.Context, req *email.MerchantEmailGatewaySetupReq) (res *email.MerchantEmailGatewaySetupRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	err = email2.SetupMerchantEmailConfig(ctx, _interface.GetMerchantId(ctx), req.GatewayName, req.Data, req.IsDefault)
	if err != nil {
		return nil, err
	}
	return &email.MerchantEmailGatewaySetupRes{}, nil
}
