package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) UpdateMemberRole(ctx context.Context, req *member.UpdateMemberRoleReq) (res *member.UpdateMemberRoleRes, err error) {
	err = member2.UpdateMemberRole(ctx, _interface.GetMerchantId(ctx), req.MemberId, req.Roles)
	if err != nil {
		return nil, err
	}
	return &member.UpdateMemberRoleRes{}, nil
}
