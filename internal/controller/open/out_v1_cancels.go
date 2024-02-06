package open

import (
	"context"
	"unibee-api/api/open/payment"
	"unibee-api/internal/logic/payment/service"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerPayment) Cancels(ctx context.Context, req *payment.CancelsReq) (res *payment.CancelsRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	overseaPay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(overseaPay != nil, "payment not found")
	err = service.PaymentGatewayCancel(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	return &payment.CancelsRes{}, nil
}
