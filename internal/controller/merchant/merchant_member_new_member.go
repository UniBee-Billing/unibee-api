package merchant

import (
	"context"
	"unibee/api/merchant/member"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"
	"unibee/utility"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	utility.Assert(!config.IsOpenSourceVersion(), "This is a premium feature")
	err = member2.AddMerchantMember(ctx, _interface.GetMerchantId(ctx), req.Email, req.FirstName, req.LastName, req.RoleIds)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
