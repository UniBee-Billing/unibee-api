package merchant

import (
	"context"
	"unibee/api/bean/detail"
	"unibee/internal/logic/jwt"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/auth"
)

func (c *ControllerAuth) OauthMembers(ctx context.Context, req *auth.OauthMembersReq) (res *auth.OauthMembersRes, err error) {
	authJsToken := jwt.GetOAuthJsTokenFromHeader(ctx)
	utility.Assert(len(authJsToken) > 0, "Auth.js JWT token is required")

	// 2. Validate Auth.js JWT token and extract user information
	authJsClaims, err := jwt.ValidateOAuthJsJWT(ctx, authJsToken)
	utility.AssertError(err, "Invalid Auth.js JWT token")

	// 3. Validate email format and match with token
	oauthMembers := query.GetMerchantMembersByAuthJsProvider(ctx, authJsClaims.Provider, authJsClaims.ProviderId)
	memberDetails := make([]*detail.MerchantMemberDetail, 0)
	for _, member := range oauthMembers {
		memberDetails = append(memberDetails, detail.ConvertMemberToDetail(ctx, member))
	}

	return &auth.OauthMembersRes{MerchantMembers: memberDetails}, nil
}
