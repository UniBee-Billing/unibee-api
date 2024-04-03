package merchant

import (
	"context"
	"unibee/api/merchant/payment"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/detail"
)

func (c *ControllerPayment) Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error) {
	return &payment.DetailRes{PaymentDetail: detail.GetPaymentDetail(ctx, _interface.GetMerchantId(ctx), req.PaymentId)}, nil
}
