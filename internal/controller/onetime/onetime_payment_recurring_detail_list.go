package onetime

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/onetime/payment"
)

func (c *ControllerPayment) RecurringDetailList(ctx context.Context, req *payment.RecurringDetailListReq) (res *payment.RecurringDetailListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
