package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	memberLogic "unibee/internal/logic/member"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) UpdateOAuth(ctx context.Context, req *member.UpdateOAuthReq) (res *member.UpdateOAuthRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Not support API")
	// 1. Get and validate Auth.js JWT token from header
	authJsToken := jwt.GetOAuthJsTokenFromHeader(ctx)
	utility.Assert(len(authJsToken) > 0, "Auth.js JWT token is required")

	// 2. Validate Auth.js JWT token and extract user information
	authJsClaims, err := jwt.ValidateOAuthJsJWT(ctx, authJsToken)
	utility.AssertError(err, "Invalid Auth.js JWT token")

	utility.Assert(len(authJsClaims.Provider) > 0, "Invalid Provider from oauth token")
	utility.Assert(len(authJsClaims.ProviderId) > 0, "Invalid ProviderId from oauth token")

	// 1. Get current merchant member from context
	currentMember := _interface.Context().Get(ctx).MerchantMember
	utility.Assert(currentMember != nil, "Merchant Member Not Found")

	// 2. Update OAuth provider for current member
	memberLogic.UpdateMemberAuthJsProvider(ctx, currentMember.Id, authJsClaims)

	return &member.UpdateOAuthRes{}, nil
}
