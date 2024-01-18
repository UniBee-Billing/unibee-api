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

	payment := query.GetPaymentByPaymentId(ctx, req.PaymentsPspReference)
	utility.Assert(payment != nil, "payment not found")
	utility.Assert(payment.Currency == req.Amount.Currency, "Currency not match the payment")
	payment.ChannelPaymentFee = req.Amount.Value
	err = service.DoChannelCapture(ctx, payment)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
