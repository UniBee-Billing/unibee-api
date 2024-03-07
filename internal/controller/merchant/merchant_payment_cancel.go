package merchant

import (
	"context"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) Cancel(ctx context.Context, req *payment.CancelReq) (res *payment.CancelRes, err error) {

	overseaPay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(overseaPay != nil, "payment not found")
	err = service.PaymentGatewayCancel(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	return &payment.CancelRes{}, nil
}
