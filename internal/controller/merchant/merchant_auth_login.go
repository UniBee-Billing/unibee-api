package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/member"
	"unibee/internal/logic/totp/client_activity"
	"unibee/utility"
)

func (c *ControllerAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	utility.Assert(req.Email != "", "Email Cannot Be Empty")
	utility.Assert(req.Password != "", "Password Cannot Be Empty")
	one, token := member.PasswordLogin(ctx, req.Email, req.Password, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity)
	utility.Assert(one.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")

	authJsToken := jwt.GetOAuthJsTokenFromHeader(ctx)
	if len(authJsToken) > 0 && len(req.Provider) > 0 && len(req.ProviderId) > 0 {
		authJsClaims, err := jwt.ValidateOAuthJsJWT(ctx, authJsToken)
		if err == nil && authJsClaims != nil && authJsClaims.ProviderId == req.ProviderId && authJsClaims.Provider == req.Provider {
			member.UpdateMemberAuthJsProvider(ctx, one.MerchantId, authJsClaims)
		} else if err != nil {
			g.Log().Errorf(ctx, "Login Connect to OAuth err:%s", err.Error())
		}
	}

	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)
	return &auth.LoginRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one), Token: token}, nil
}
