package user

import (
	"context"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/auth"
	"unibee-api/utility"

	"unibee-api/api/user/profile"
)

func (c *ControllerProfile) PasswordReset(ctx context.Context, req *profile.PasswordResetReq) (res *profile.PasswordResetRes, err error) {
	//User 检查
	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
	utility.Assert(_interface.BizCtx().Get(ctx).User.Id > 0, "userId invalid")
	auth.ChangeUserPassword(ctx, _interface.GetMerchantId(ctx), _interface.BizCtx().Get(ctx).User.Email, req.OldPassword, req.NewPassword)
	return &profile.PasswordResetRes{}, nil
}
