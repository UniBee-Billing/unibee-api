package user

import (
	"context"
	"go-oversea-pay/api/user/auth"
	_interface "go-oversea-pay/internal/interface"
	auth2 "go-oversea-pay/internal/logic/auth"
	"go-oversea-pay/utility"
	// entity "go-oversea-pay/internal/model/entity/oversea_pay"
	// "go-oversea-pay/internal/query"
	// "github.com/gogf/gf/v2/errors/gcode"
	// "github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerAuth) Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error) {
	utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "User Not Found")
	utility.Assert(len(_interface.BizCtx().Get(ctx).User.Token) > 0, "Token Not Found")
	auth2.DelAuthToken(ctx, _interface.BizCtx().Get(ctx).User.Token)
	return &auth.LogoutRes{}, nil
}
