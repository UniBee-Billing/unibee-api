package open

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/open/payment"
)

func (c *ControllerPayment) ListRecurringDetails(ctx context.Context, req *payment.ListRecurringDetailsReq) (res *payment.ListRecurringDetailsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
