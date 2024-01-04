package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/profile"
)

func (c *ControllerProfile) Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
