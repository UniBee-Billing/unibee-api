package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/method"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) MethodDelete(ctx context.Context, req *payment.MethodDeleteReq) (res *payment.MethodDeleteRes, err error) {
	err = method.DeletePaymentMethod(ctx, _interface.GetMerchantId(ctx), _interface.Context().Get(ctx).User.Id, req.GatewayId, req.PaymentMethodId)
	if err != nil {
		return nil, err
	}
	return &payment.MethodDeleteRes{}, nil
}
