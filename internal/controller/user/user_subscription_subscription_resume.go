package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) Resume(ctx context.Context, req *subscription.ResumeReq) (res *subscription.ResumeRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
