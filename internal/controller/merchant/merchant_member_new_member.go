package merchant

import (
	"context"
	"unibee/api/merchant/member"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"
	"unibee/internal/logic/middleware"
	"unibee/utility"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	utility.Assert(middleware.IsPremiumVersion(ctx, _interface.GetMerchantId(ctx)), "Feature member need premium license, contact us directly if needed")
	err = member2.AddMerchantMember(ctx, _interface.GetMerchantId(ctx), req.Email, req.FirstName, req.LastName, req.RoleIds)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
