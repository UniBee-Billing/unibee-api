package merchant

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/user"
)

func (c *ControllerUser) Logout(ctx context.Context, req *user.LogoutReq) (res *user.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")
	auth.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).MerchantUser.Token)
	return &user.LogoutRes{}, nil
}
