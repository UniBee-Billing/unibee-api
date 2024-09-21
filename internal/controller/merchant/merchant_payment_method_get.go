package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/method"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) MethodGet(ctx context.Context, req *payment.MethodGetReq) (res *payment.MethodGetRes, err error) {
	return &payment.MethodGetRes{Method: method.QueryPaymentMethod(ctx, _interface.GetMerchantId(ctx), req.UserId, req.GatewayId, req.PaymentMethodId)}, nil
}
