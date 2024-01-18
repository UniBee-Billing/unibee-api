package open

import (
	"context"
	"go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/utility"
)

func (c *ControllerPayment) Refunds(ctx context.Context, req *payment.RefundsReq) (res *payment.RefundsRes, err error) {
	utility.Assert(req != nil, "req should not be nil")
	utility.Assert(len(req.PaymentsPspReference) > 0, "PaymentsPspReference should not be nil")
	utility.Assert(req.Amount != nil, "Amount should not be nil")
	utility.Assert(req.Amount.Value > 0, "refund value should > 0")
	utility.Assert(len(req.Amount.Currency) > 0, "refund currency should not be nil")
	currencyNumberCheck(req.Amount)
	//参数有效性校验 todo mark
	openApiConfig, _ := merchantCheck(ctx, req.MerchantId)

	resp, err := service.DoChannelRefund(ctx, consts.PAYMENT_BIZ_TYPE_ORDER, req, int64(openApiConfig.Id))
	if err != nil {
		return nil, err
	}
	res = &payment.RefundsRes{
		Status:              "SentForRefund",
		PspReference:        resp.RefundId,
		Reference:           req.Reference,
		PaymentPspReference: resp.PaymentId,
	}
	return res, nil
}
