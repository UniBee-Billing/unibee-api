package user

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/user/auth"
	_interface "unibee/internal/interface"
	auth2 "unibee/internal/logic/auth"
	"unibee/utility"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.Password != "", "Password Cannot Be Empty")
	one, token := auth2.PasswordLogin(ctx, _interface.GetMerchantId(ctx), req.Email, req.Password)
	g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
	return &auth.LoginRes{User: bean.SimplifyUserAccount(one), Token: token}, nil
}
