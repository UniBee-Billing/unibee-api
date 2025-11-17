package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/detail"

	"unibee/api/user/payment"
)

func (c *ControllerPayment) Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error) {
	one := detail.GetPaymentDetail(ctx, _interface.GetMerchantId(ctx), req.PaymentId)
	return &payment.DetailRes{PaymentDetail: one, PaymentStatus: one.Payment.Status}, nil
}
