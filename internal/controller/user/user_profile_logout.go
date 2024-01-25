package user

import (
	"context"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/user/profile"
)

func (c *ControllerProfile) Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).User.Token) > 0, "Token Not Found")
	auth.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).User.Token)
	return &profile.LogoutRes{}, nil
}
