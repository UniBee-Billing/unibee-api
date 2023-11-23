package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) PaymentMethods(ctx context.Context, req *v1.PaymentMethodsReq) (res *v1.PaymentMethodsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
