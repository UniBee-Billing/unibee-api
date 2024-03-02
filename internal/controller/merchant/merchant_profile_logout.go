package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/jwt"
	"unibee/utility"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantMember.Token) > 0, "Merchant Token Not Found")
	jwt.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).MerchantMember.Token)
	return &profile.LogoutRes{}, nil
}
