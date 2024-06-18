package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"
	"unibee/internal/logic/role"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error) {
	if req.RoleId > 0 {
		list, total := role.GetMemberListByRole(ctx, _interface.GetMerchantId(ctx), req.RoleId)
		return &member.ListRes{MerchantMembers: list, Total: total}, nil
	} else {
		list, total := member2.MerchantMemberList(ctx, _interface.GetMerchantId(ctx))
		return &member.ListRes{MerchantMembers: list, Total: total}, nil
	}

}
