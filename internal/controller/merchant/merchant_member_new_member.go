package merchant

import (
	"context"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error) {
	err = member2.AddMerchantMember(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &member.NewMemberRes{}, nil
}
