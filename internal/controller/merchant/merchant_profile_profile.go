package merchant

import (
	"context"

	"go-oversea-pay/api/merchant/profile"
)

func (c *ControllerProfile) Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error) {
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
	// to be implemented later
	return &profile.ProfileRes{}, nil
}
