package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) Captures(ctx context.Context, req *v1.CapturesReq) (res *v1.CapturesRes, err error) {
	panic("capture panic error moke")
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
