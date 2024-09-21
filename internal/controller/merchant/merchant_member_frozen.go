package merchant

import (
	"context"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Frozen(ctx context.Context, req *member.FrozenReq) (res *member.FrozenRes, err error) {
	member2.FrozenMember(ctx, req.MemberId)
	return &member.FrozenRes{}, nil
}
