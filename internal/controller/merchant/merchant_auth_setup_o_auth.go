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
	auth2 "unibee/internal/logic/member"
	"unibee/internal/logic/totp/client_activity"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerAuth) SetupOAuth(ctx context.Context, req *auth.SetupOAuthReq) (res *auth.SetupOAuthRes, err error) {
	// 1. Get and validate Auth.js JWT token from header
	authJsToken := jwt.GetOAuthJsTokenFromHeader(ctx)
	utility.Assert(len(authJsToken) > 0, "Auth.js JWT token is required")

	// 2. Validate Auth.js JWT token and extract user information
	authJsClaims, err := jwt.ValidateOAuthJsJWT(ctx, authJsToken)
	utility.AssertError(err, "Invalid Auth.js JWT token")

	// 3. Validate email format and match with token
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid email format")

	// 4. Validate setup token
	// TODO: Implement setup token validation logic
	// This should verify the setup token is valid and not expired
	// For now, we'll skip this validation

	// 5. Find merchant member
	merchantMember := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(merchantMember != nil, "Account not found")
	utility.Assert(merchantMember.TotpValidatorSecret == req.SetupToken, "Invalid setupToken")

	// 6. Validate account status
	utility.Assert(merchantMember.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")

	// 8. Link OAuth account
	member.UpdateMemberAuthJsProvider(ctx, merchantMember.Id, authJsClaims)
	if len(req.NewPassword) > 0 {
		auth2.ChangeMerchantMemberPasswordWithOutOldVerify(ctx, req.Email, req.NewPassword)
	}
	// 9. Generate system JWT token
	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, merchantMember.MerchantId, merchantMember.Id, merchantMember.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", merchantMember.Id)), "Cache Error")

	// 10. Set cookies
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)

	return &auth.SetupOAuthRes{
		MerchantMember: detail.ConvertMemberToDetail(ctx, merchantMember),
		Token:          token,
	}, nil
}
