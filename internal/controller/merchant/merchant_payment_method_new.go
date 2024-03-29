package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/method"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) MethodNew(ctx context.Context, req *payment.MethodNewReq) (res *payment.MethodNewRes, err error) {
	return &payment.MethodNewRes{Method: method.NewPaymentMethod(ctx, &method.NewPaymentMethodInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     req.UserId,
		GatewayId:  req.GatewayId,
		Type:       req.Type,
		Data:       req.Data,
	})}, nil
}
