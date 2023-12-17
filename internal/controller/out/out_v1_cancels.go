package out

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/out/v1"
	"go-oversea-pay/internal/logic/payment/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerV1) Cancels(ctx context.Context, req *v1.CancelsReq) (res *v1.CancelsRes, err error) {
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantId)

	overseaPay := query.GetOverseaPayByMerchantOrderNo(ctx, req.PaymentsPspReference)
	utility.Assert(overseaPay != nil, "payment not found")
	err = service.DoChannelCancel(ctx, overseaPay)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
