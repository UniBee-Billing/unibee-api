package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) Suspend(ctx context.Context, req *subscription.SuspendReq) (res *subscription.SuspendRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
