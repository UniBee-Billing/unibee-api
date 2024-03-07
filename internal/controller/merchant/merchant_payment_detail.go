package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/query"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error) {
	return &payment.DetailRes{PaymentDetail: query.GetPaymentDetail(ctx, _interface.GetMerchantId(ctx), req.PaymentId)}, nil
}
