package user

import (
	"context"
	"fmt"
	"unibee/api/bean"
	"unibee/api/user/auth"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"

	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func (c *ControllerAuth) LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error) {
	verificationCode, err := g.Redis().Get(ctx, req.Email+"-Verify")
	utility.AssertError(err, "Server Error")
	utility.Assert(verificationCode != nil, "code expired")
	utility.Assert((verificationCode.String()) == req.VerificationCode, "code not match")
	var one *entity.UserAccount
	one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	utility.Assert(one != nil, "Login Failed")
	token, err := jwt.CreatePortalToken(jwt.TOKENTYPEUSER, one.MerchantId, one.Id, req.Email, one.Language)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", one.Id)), "Cache Error")
	g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	return &auth.LoginOtpVerifyRes{User: bean.SimplifyUserAccount(one), Token: token}, nil
}
