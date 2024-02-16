package merchant

import (
	"context"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/auth"
	"unibee-api/utility"

	"unibee-api/api/merchant/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")
	auth.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).MerchantUser.Token)
	return &profile.LogoutRes{}, nil
}
