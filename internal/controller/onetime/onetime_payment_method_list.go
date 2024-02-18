package onetime

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/onetime/payment"
)

func (c *ControllerPayment) PaymentMethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
