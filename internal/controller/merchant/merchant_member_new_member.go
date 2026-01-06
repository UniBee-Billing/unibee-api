package merchant

import (
	"context"
	"fmt"
	"unibee/api/merchant/member"
	_interface "unibee/internal/interface/context"
	member2 "unibee/internal/logic/member"
	"unibee/internal/logic/middleware/license"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	utility.Assert(license.IsPremiumVersion(ctx, _interface.GetMerchantId(ctx)), "Feature member need premium license, contact us directly if needed")
	maxMemberCount := license.GetMerchantMemberLimit(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(maxMemberCount <= 0 || maxMemberCount > query.GetMerchantMemberCount(ctx, _interface.GetMerchantId(ctx)), fmt.Sprintf("You have reached max members limit: %d, please upgrade plan", maxMemberCount))
	err = member2.AddMerchantMember(ctx, _interface.GetMerchantId(ctx), req.Email, req.FirstName, req.LastName, req.RoleIds, req.ReturnUrl)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
