package onetime

import (
	"context"
	"unibee/api/onetime/payment"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/payment/service"
	"unibee/utility"
)

func (c *ControllerPayment) NewPaymentRefund(ctx context.Context, req *payment.NewPaymentRefundReq) (res *payment.NewPaymentRefundRes, err error) {
	utility.Assert(req != nil, "req should not be nil")
	utility.Assert(len(req.PaymentId) > 0, "PaymentId should not be nil")
	utility.Assert(req.RefundAmount > 0, "refund value should > 0")
	utility.Assert(len(req.Currency) > 0, "refund currency should not be nil")
	currencyNumberCheck(req.RefundAmount, req.Currency)
	openApiConfig, _ := merchantCheck(ctx, _interface.GetMerchantId(ctx))

	resp, err := service.GatewayPaymentRefundCreate(ctx, consts.BizTypeOneTime, req, int64(openApiConfig.Id))
	if err != nil {
		return nil, err
	}
	res = &payment.NewPaymentRefundRes{
		Status:           "SentForRefund",
		RefundId:         resp.RefundId,
		MerchantRefundId: req.ExternalRefundId,
		PaymentId:        resp.PaymentId,
	}
	return res, nil
}
