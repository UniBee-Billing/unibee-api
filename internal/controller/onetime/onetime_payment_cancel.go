package onetime

import (
	"context"
	"unibee-api/api/onetime/payment"
	"unibee-api/internal/logic/payment/service"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerPayment) Cancel(ctx context.Context, req *payment.CancelReq) (res *payment.CancelRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	overseaPay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(overseaPay != nil, "payment not found")
	err = service.PaymentGatewayCancel(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	return &payment.CancelRes{}, nil
}
