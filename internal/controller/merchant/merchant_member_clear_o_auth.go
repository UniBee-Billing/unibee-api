package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	memberLogic "unibee/internal/logic/member"
	"unibee/utility"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) ClearOAuth(ctx context.Context, req *member.ClearOAuthReq) (res *member.ClearOAuthRes, err error) {
	// Validate current user permissions
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "Not support API")

	// Validate request parameters
	utility.Assert(len(req.Provider) > 0, "Provider is required")

	// Get current user ID
	memberId := _interface.Context().Get(ctx).MerchantMember.Id

	// Call logic layer function to remove OAuth binding
	memberLogic.ClearMemberAuthJsProvider(ctx, memberId, req.Provider)

	return &member.ClearOAuthRes{}, nil
}
