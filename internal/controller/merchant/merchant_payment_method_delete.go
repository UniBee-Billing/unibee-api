package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/method"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) MethodDelete(ctx context.Context, req *payment.MethodDeleteReq) (res *payment.MethodDeleteRes, err error) {
	err = method.DeletePaymentMethod(ctx, _interface.GetMerchantId(ctx), req.UserId, req.GatewayId, req.PaymentMethodId)
	if err != nil {
		return nil, err
	}
	return &payment.MethodDeleteRes{}, nil
}
