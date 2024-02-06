package open

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/open/payment"
)

func (c *ControllerPayment) PaymentMethods(ctx context.Context, req *payment.PaymentMethodsReq) (res *payment.PaymentMethodsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
