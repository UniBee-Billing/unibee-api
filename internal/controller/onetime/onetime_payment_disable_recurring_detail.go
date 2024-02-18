package onetime

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/onetime/payment"
)

func (c *ControllerPayment) DisableRecurringDetail(ctx context.Context, req *payment.DisableRecurringDetailReq) (res *payment.DisableRecurringDetailRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
