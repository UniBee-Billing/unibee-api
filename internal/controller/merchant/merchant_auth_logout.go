package merchant

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
	auth2 "go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/auth"
)

func (c *ControllerAuth) Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")
	auth2.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).MerchantUser.Token)
	return &auth.LogoutRes{}, nil
}
