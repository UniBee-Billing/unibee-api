package merchant

import (
	"context"
	"unibee/api/merchant/payment"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/payment/detail"
)

func (c *ControllerPayment) RefundDetail(ctx context.Context, req *payment.RefundDetailReq) (res *payment.RefundDetailRes, err error) {
	return &payment.RefundDetailRes{RefundDetail: detail.GetRefundDetail(ctx, _interface.GetMerchantId(ctx), req.RefundId)}, nil
}
