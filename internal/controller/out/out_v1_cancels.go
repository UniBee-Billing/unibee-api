package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) Cancels(ctx context.Context, req *v1.CancelsReq) (res *v1.CancelsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
