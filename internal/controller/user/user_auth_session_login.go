package user

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/internal/logic/jwt"
	auth2 "unibee/internal/logic/session"
	"unibee/utility"

	"unibee/api/user/auth"
)

func (c *ControllerAuth) SessionLogin(ctx context.Context, req *auth.SessionLoginReq) (res *auth.SessionLoginRes, err error) {
	one, returnUrl := auth2.UserSessionTransfer(ctx, req.Session)
	utility.Assert(one != nil, "Login Failed")

	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, one.Email)
	fmt.Println("session logged-in, save email/id in token: ", one.Email, "/", one.Id)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
	return &auth.SessionLoginRes{User: bean.SimplifyUserAccount(one), Token: token, ReturnUrl: returnUrl}, nil
}
