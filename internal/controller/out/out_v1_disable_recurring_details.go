package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) DisableRecurringDetails(ctx context.Context, req *v1.DisableRecurringDetailsReq) (res *v1.DisableRecurringDetailsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
