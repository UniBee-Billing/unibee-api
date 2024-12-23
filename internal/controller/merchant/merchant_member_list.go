package merchant

import (
	"context"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	member2 "unibee/internal/logic/member"
	"unibee/internal/logic/role"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error) {
	if req.RoleIds != nil && len(req.RoleIds) > 0 {
		if req.Page > 1 {
			return &member.ListRes{MerchantMembers: make([]*detail.MerchantMemberDetail, 0), Total: 0}, nil
		}
		list, total := role.GetMemberListByRoleIds(ctx, _interface.GetMerchantId(ctx), req.RoleIds, req.Page, req.Count)
		return &member.ListRes{MerchantMembers: list, Total: total}, nil
	} else {
		list, total := member2.MerchantMemberList(ctx, _interface.GetMerchantId(ctx), req.Page, req.Count)
		return &member.ListRes{MerchantMembers: list, Total: total}, nil
	}

}
