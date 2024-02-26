package onetime

import (
	"context"
	"unibee/api/onetime/payment"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPayment) Capture(ctx context.Context, req *payment.CaptureReq) (res *payment.CaptureRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, _interface.GetMerchantId(ctx))

	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	utility.Assert(one.Currency == req.Amount.Currency, "Currency not match the payment")
	one.PaymentAmount = req.Amount.Amount
	err = service.PaymentGatewayCapture(ctx, one)
	if err != nil {
		return nil, err
	}
	return &payment.CaptureRes{}, nil
}
