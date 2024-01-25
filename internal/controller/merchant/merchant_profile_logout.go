package merchant

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")
	auth.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).MerchantUser.Token)
	return &profile.LogoutRes{}, nil
}
