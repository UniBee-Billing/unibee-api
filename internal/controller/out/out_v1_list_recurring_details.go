package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) ListRecurringDetails(ctx context.Context, req *v1.ListRecurringDetailsReq) (res *v1.ListRecurringDetailsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
