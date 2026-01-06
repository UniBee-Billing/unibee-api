package merchant

import (
	"context"
	"fmt"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/member"
	"unibee/internal/logic/totp"
	"unibee/internal/logic/totp/client_activity"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) LoginOAuth(ctx context.Context, req *auth.LoginOAuthReq) (res *auth.LoginOAuthRes, err error) {
	// 1. Get and validate Auth.js JWT token from header
	authJsToken := jwt.GetOAuthJsTokenFromHeader(ctx)
	utility.Assert(len(authJsToken) > 0, "Auth.js JWT token is required")

	// 2. Validate Auth.js JWT token and extract user information
	authJsClaims, err := jwt.ValidateOAuthJsJWT(ctx, authJsToken)
	utility.AssertError(err, "Invalid Auth.js JWT token")

	// 3. Validate email format and match with token
	utility.Assert(utility.IsEmailValid(req.Email), "Invalid email format")
	utility.Assert(len(authJsClaims.Provider) > 0, "Invalid Provider from oauth token")
	utility.Assert(len(authJsClaims.ProviderId) > 0, "Invalid ProviderId from oauth token")
	merchantMember := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(merchantMember != nil, "Account not found. Please register first.")

	// 4. First try to find user by OAuth information
	oauthMembers := query.GetMerchantMembersByAuthJsProvider(ctx, authJsClaims.Provider, authJsClaims.ProviderId)
	oauthConnected := false
	for _, one := range oauthMembers {
		if one.Id == merchantMember.Id {
			oauthConnected = true
		}
	}

	if merchantMember.Email == authJsClaims.Email {
		member.UpdateMemberAuthJsProvider(ctx, merchantMember.Id, authJsClaims)
		oauthConnected = true
	}
	utility.Assert(oauthConnected == true, "OAuth Account not connected")

	// 7. Validate account status
	utility.Assert(merchantMember.Status != 2, "Your account has been suspended. Please contact billing admin for further assistance.")

	// 8. 2FA validation (if enabled) - same logic as password login
	if merchantMember.TotpValidatorType > 0 && len(merchantMember.TotpValidatorSecret) > 0 && !totp.IsClientIdentityValid(ctx, merchantMember.Email, _interface.Context().Get(ctx).ClientIdentity) {
		// need totp validate, not needed for every login
		utility.Assert(len(req.TotpCode) > 0, "2FA Code Needed")
		utility.Assert(totp.ValidateMerchantMemberTotp(ctx, merchantMember.TotpValidatorType, merchantMember.Email, merchantMember.TotpValidatorSecret, req.TotpCode, _interface.Context().Get(ctx).ClientIdentity), "Invalid 2FA Code")
	}

	// 9. Generate system JWT token
	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, merchantMember.MerchantId, merchantMember.Id, merchantMember.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", merchantMember.Id)), "Cache Error")

	// 10. Set cookies
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)

	return &auth.LoginOAuthRes{
		MerchantMember: detail.ConvertMemberToDetail(ctx, merchantMember),
		Token:          token,
	}, nil
}
