package merchant

import (
	"context"
	"unibee/api/merchant/member"
	_interface "unibee/internal/interface/context"
	member2 "unibee/internal/logic/member"
)

func (c *ControllerMember) List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error) {
	list, total := member2.MerchantMemberList(ctx, &member2.MemberListReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		SearchKey:       req.SearchKey,
		Email:           req.Email,
		RoleIds:         req.RoleIds,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	return &member.ListRes{MerchantMembers: list, Total: total}, nil
}
