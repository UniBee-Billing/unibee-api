package out

import (
	"context"
	"go-oversea-pay/api/out/v1"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/service"
)

func (c *ControllerV1) Refunds(ctx context.Context, req *v1.RefundsReq) (res *v1.RefundsRes, err error) {
	currencyNumberCheck(req.Amount)
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	// openApiId todo mark
	resp, err := service.DoChannelRefund(ctx, consts.PAYMENT_BIZ_TYPE_ORDER, req, 0)
	if err != nil {
		return nil, err
	}
	res = &v1.RefundsRes{
		Status:              "SentForRefund",
		PspReference:        resp.OutRefundNo,
		Reference:           req.Reference,
		PaymentPspReference: resp.OutTradeNo,
	}
	return res, nil
}
