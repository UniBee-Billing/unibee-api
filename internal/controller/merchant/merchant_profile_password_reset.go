package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/utility"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) PasswordReset(ctx context.Context, req *profile.PasswordResetReq) (res *profile.PasswordResetRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "Merchant User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).MerchantUser.Token) > 0, "Merchant Token Not Found")

	auth.ChangeMerchantUserPassword(ctx, _interface.GetMerchantId(ctx), _interface.BizCtx().Get(ctx).MerchantUser.Email, req.OldPassword, req.NewPassword)
	return &profile.PasswordResetRes{}, nil
}
