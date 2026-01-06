package merchant

import (
	"context"
	"fmt"
	"strings"
	"unibee/api/bean/detail"
	"unibee/api/merchant/auth"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/member"
	"unibee/internal/logic/merchant"
	"unibee/internal/logic/middleware/license"
	"unibee/internal/logic/totp/client_activity"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerAuth) RegisterOAuth(ctx context.Context, req *auth.RegisterOAuthReq) (res *auth.RegisterOAuthRes, err error) {
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

	// 4. Check if email already exists
	existingMember := query.GetMerchantMemberByEmail(ctx, req.Email)
	utility.Assert(existingMember == nil, fmt.Sprintf("Merchant With Email (%s) already exists", req.Email))

	// 6. Check merchant limit and license (same logic as regular register)
	list := query.GetActiveMerchantList(ctx)
	if len(list) > 2 {
		utility.Assert(config.GetConfigInstance().Mode == "cloud", "Register multi merchants should contain valid mode")
		var containPremiumMerchant = false
		for _, one := range list {
			if license.IsPremiumVersion(ctx, one.Id) {
				containPremiumMerchant = true
				break
			}
		}
		utility.Assert(containPremiumMerchant, "Feature register multi merchants need premium license, contact us directly if needed")
	}

	// 7. Create merchant and member using the same logic as regular register
	// Use Auth.js data if available, otherwise use request data
	firstName := req.FirstName
	lastName := req.LastName
	if firstName == "" && authJsClaims.Name != "" {
		// Try to extract first and last name from Auth.js name
		names := strings.Split(authJsClaims.Name, " ")
		if len(names) > 0 {
			firstName = names[0]
		}
		if len(names) > 1 {
			lastName = strings.Join(names[1:], " ")
		}
	}

	createMerchantReq := &merchant.CreateMerchantInternalReq{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       req.Email,
		Password:    req.Password,
		Phone:       req.Phone,
		UserName:    req.UserName,
		CountryCode: req.CountryCode,
		CountryName: req.CountryName,
		CompanyName: req.CompanyName,
		Metadata:    req.Metadata,
	}

	_, merchantMember, err := merchant.CreateMerchant(ctx, createMerchantReq)
	utility.AssertError(err, "Failed to create merchant")

	// 8. Link OAuth account
	member.UpdateMemberAuthJsProvider(ctx, merchantMember.Id, authJsClaims)

	// 9. Generate system JWT token
	token, err := jwt.CreateMemberPortalToken(ctx, jwt.TOKENTYPEMERCHANTMember, merchantMember.MerchantId, merchantMember.Id, merchantMember.Email)
	utility.AssertError(err, "Server Error")
	utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("MerchantMember#%d", merchantMember.Id)), "Cache Error")

	// 10. Set cookies
	g.RequestFromCtx(ctx).Cookie.Set(jwt.MERCHANT_TYPE_TOKEN_COOKIE_KEY, token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	client_activity.UpdateClientIdentityLoginTime(ctx, _interface.Context().Get(ctx).ClientIdentity)

	return &auth.RegisterOAuthRes{
		MerchantMember: detail.ConvertMemberToDetail(ctx, merchantMember),
		Token:          token,
	}, nil
}
