package out

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) PaymentDetails(ctx context.Context, req *v1.PaymentDetailsReq) (res *v1.PaymentDetailsRes, err error) {
	panic("capture panic error moke")
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
