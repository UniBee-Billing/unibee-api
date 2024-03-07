package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) RefundList(ctx context.Context, req *payment.RefundListReq) (res *payment.RefundListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
