package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	err = member2.AddMerchantMember(ctx, _interface.GetMerchantId(ctx), req.Email, req.FirstName, req.LastName, req.Roles)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
