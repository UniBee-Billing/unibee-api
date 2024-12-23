package user

import (
	"context"
	"unibee/api/user/payment"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/method"
)

func (c *ControllerPayment) MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error) {
	return &payment.MethodListRes{MethodList: method.QueryPaymentMethodList(ctx, &method.PaymentMethodListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     _interface.Context().Get(ctx).User.Id,
		GatewayId:  req.GatewayId,
		PaymentId:  req.PaymentId,
	})}, nil
}
