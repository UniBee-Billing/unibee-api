package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/method"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) MethodGet(ctx context.Context, req *payment.MethodGetReq) (res *payment.MethodGetRes, err error) {
	return &payment.MethodGetRes{Method: method.QueryPaymentMethod(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, req.GatewayId, req.PaymentMethodId)}, nil
}
