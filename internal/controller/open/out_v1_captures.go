package open

import (
	"context"
	"go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerPayment) Captures(ctx context.Context, req *payment.CapturesReq) (res *payment.CapturesRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	utility.Assert(one.Currency == req.Amount.Currency, "Currency not match the payment")
	one.ChannelPaymentFee = req.Amount.Value
	err = service.DoChannelCapture(ctx, one)
	if err != nil {
		return nil, err
	}
	return &payment.CapturesRes{}, nil
}
