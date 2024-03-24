package merchant

import (
	"context"
	"strings"
	"unibee/internal/consts"
	"unibee/internal/logic/payment/service"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) NewPaymentRefund(ctx context.Context, req *payment.NewPaymentRefundReq) (res *payment.NewPaymentRefundRes, err error) {
	utility.Assert(req != nil, "req should not be nil")
	utility.Assert(len(req.PaymentId) > 0, "PaymentId should not be nil")
	utility.Assert(req.RefundAmount > 0, "refund value should > 0")
	req.Currency = strings.ToUpper(req.Currency)
	utility.Assert(len(req.Currency) > 0, "refund currency should not be nil")
	currencyNumberCheck(req.RefundAmount, req.Currency)

	resp, err := service.GatewayPaymentRefundCreate(ctx, &service.NewPaymentRefundInternalReq{
		PaymentId:        req.PaymentId,
		ExternalRefundId: req.ExternalRefundId,
		RefundAmount:     req.RefundAmount,
		Currency:         req.Currency,
		Reason:           req.Reason,
		Metadata:         req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	res = &payment.NewPaymentRefundRes{
		Status:           consts.RefundCreated,
		RefundId:         resp.RefundId,
		ExternalRefundId: req.ExternalRefundId,
		PaymentId:        resp.PaymentId,
	}
	return res, nil
}
