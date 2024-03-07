package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) RefundDetail(ctx context.Context, req *payment.RefundDetailReq) (res *payment.RefundDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
