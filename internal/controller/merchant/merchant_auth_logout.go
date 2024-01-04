package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/auth"
)

func (c *ControllerAuth) Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
