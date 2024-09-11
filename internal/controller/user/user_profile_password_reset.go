package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/user"
	"unibee/utility"

	"unibee/api/user/profile"
)

func (c *ControllerProfile) PasswordReset(ctx context.Context, req *profile.PasswordResetReq) (res *profile.PasswordResetRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
	utility.Assert(_interface.Context().Get(ctx).User.Id > 0, "userId invalid")
	user.ChangeUserPassword(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Email, req.OldPassword, req.NewPassword)
	return &profile.PasswordResetRes{}, nil
}
