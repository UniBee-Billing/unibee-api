package webhook

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/webhook/setup"
)

func (c *ControllerSetup) New(ctx context.Context, req *setup.NewReq) (res *setup.NewRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
