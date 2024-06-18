package merchant

import (
	"context"
	member2 "unibee/internal/logic/member"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Release(ctx context.Context, req *member.ReleaseReq) (res *member.ReleaseRes, err error) {
	member2.ReleaseMember(ctx, req.MemberId)
	return &member.ReleaseRes{}, nil
}
