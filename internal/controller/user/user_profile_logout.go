package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/jwt"
	"unibee/utility"

	"unibee/api/user/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).User != nil, "User Not Found")
	utility.Assert(len(_interface.Context().Get(ctx).User.Token) > 0, "Token Not Found")
	jwt.DelAuthToken(ctx, _interface.Context().Get(ctx).User.Token)
	return &profile.LogoutRes{}, nil
}
