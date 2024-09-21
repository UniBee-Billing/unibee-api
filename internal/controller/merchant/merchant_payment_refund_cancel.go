package merchant

import (
	"context"
	"unibee/internal/logic/payment/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) RefundCancel(ctx context.Context, req *payment.RefundCancelReq) (res *payment.RefundCancelRes, err error) {
	one := query.GetRefundByRefundId(ctx, req.RefundId)
	utility.Assert(one != nil, "refund not found")
	err = service.PaymentRefundGatewayCancel(ctx, one)
	if err != nil {
		return nil, err
	}
	return &payment.RefundCancelRes{}, nil
}
