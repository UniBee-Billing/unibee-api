package merchant

import (
	"context"

	"unibee/api/merchant/member"
)

func (c *ControllerMember) Profile(ctx context.Context, req *member.ProfileReq) (res *member.ProfileRes, err error) {
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
	// to be implemented later
	return &member.ProfileRes{}, nil
}
