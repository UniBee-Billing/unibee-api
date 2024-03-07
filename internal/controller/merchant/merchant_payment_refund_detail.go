package merchant

import (
	"context"
	"unibee/api/merchant/payment"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
)

func (c *ControllerPayment) RefundDetail(ctx context.Context, req *payment.RefundDetailReq) (res *payment.RefundDetailRes, err error) {
	return &payment.RefundDetailRes{RefundDetail: query.GetRefundDetail(ctx, _interface.GetMerchantId(ctx), req.RefundId)}, nil
}
