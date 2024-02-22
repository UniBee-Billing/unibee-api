package user

import (
	"context"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/jwt"
	"unibee-api/utility"

	"unibee-api/api/user/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).User.Token) > 0, "Token Not Found")
	jwt.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).User.Token)
	return &profile.LogoutRes{}, nil
}
