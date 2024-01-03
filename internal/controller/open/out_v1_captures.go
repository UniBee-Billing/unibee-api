package open

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerPayment) Captures(ctx context.Context, req *payment.CapturesReq) (res *payment.CapturesRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	overseaPay := query.GetOverseaPayByMerchantOrderNo(ctx, req.PaymentsPspReference)
	utility.Assert(overseaPay != nil, "payment not found")
	utility.Assert(overseaPay.Currency == req.Amount.Currency, "Currency not match the payment")
	overseaPay.BuyerPayFee = req.Amount.Value
	err = service.DoChannelCapture(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
