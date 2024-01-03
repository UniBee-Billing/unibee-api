package open

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/api/open/payment"
)

func (c *ControllerPayment) PaymentDetails(ctx context.Context, req *payment.PaymentDetailsReq) (res *payment.PaymentDetailsRes, err error) {
	panic("capture panic error moke")
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
