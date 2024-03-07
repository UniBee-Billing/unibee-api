package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) List(ctx context.Context, req *payment.ListReq) (res *payment.ListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
