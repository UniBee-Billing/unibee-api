package open

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerPayment) Cancels(ctx context.Context, req *payment.CancelsReq) (res *payment.CancelsRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	overseaPay := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(overseaPay != nil, "payment not found")
	err = service.DoChannelCancel(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
