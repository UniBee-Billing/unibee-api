package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/member"
	"unibee/internal/logic/totp/client_activity"
	"unibee/utility"
)

func (c *ControllerAuth) Session(ctx context.Context, req *auth.SessionReq) (res *auth.SessionRes, err error) {
	one, returnUrl := member.SessionTransfer(ctx, req.Session)

	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, one.MerchantId, one.Id, one.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", one.Id)), "Cache Error")
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)
	return &auth.SessionRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one), Token: token, ReturnUrl: returnUrl}, nil
}
