package merchant

import (
	"context"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) Cancel(ctx context.Context, req *payment.CancelReq) (res *payment.CancelRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	err = service.PaymentGatewayCancel(ctx, one)
	if err != nil {
		return nil, err
	}
	return &payment.CancelRes{}, nil
}
