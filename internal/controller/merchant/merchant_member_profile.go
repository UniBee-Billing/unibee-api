package merchant

import (
	"context"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Profile(ctx context.Context, req *member.ProfileReq) (res *member.ProfileRes, err error) {
	return &member.ProfileRes{}, nil
}
