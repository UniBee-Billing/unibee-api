package merchant

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/payment"
)

func (c *ControllerPayment) EventList(ctx context.Context, req *payment.EventListReq) (res *payment.EventListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
