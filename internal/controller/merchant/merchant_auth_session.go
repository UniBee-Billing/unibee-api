package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/member"
	"unibee/utility"
)

func (c *ControllerAuth) Session(ctx context.Context, req *auth.SessionReq) (res *auth.SessionRes, err error) {
	one, returnUrl := member.SessionTransfer(ctx, req.Session)

	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, one.MerchantId, one.Id, one.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", one.Id)), "Cache Error")
	jwt.AppendRequestCookieWithToken(ctx, token)
	return &auth.SessionRes{MerchantMember: detail.ConvertMemberToDetail(ctx, one), Token: token, ReturnUrl: returnUrl}, nil
}
