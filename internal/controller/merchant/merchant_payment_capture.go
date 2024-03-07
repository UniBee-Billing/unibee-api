package merchant

import (
	"context"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) Capture(ctx context.Context, req *payment.CaptureReq) (res *payment.CaptureRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	utility.Assert(one.Currency == req.Currency, "Currency not match the payment")
	one.PaymentAmount = req.CaptureAmount
	err = service.PaymentGatewayCapture(ctx, one)
	if err != nil {
		return nil, err
	}
	return &payment.CaptureRes{}, nil
}
