package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
