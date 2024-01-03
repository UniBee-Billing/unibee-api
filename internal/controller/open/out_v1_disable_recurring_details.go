package open

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/api/open/payment"
)

func (c *ControllerPayment) DisableRecurringDetails(ctx context.Context, req *payment.DisableRecurringDetailsReq) (res *payment.DisableRecurringDetailsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
