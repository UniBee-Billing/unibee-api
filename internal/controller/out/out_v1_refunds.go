package out

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/out/v1"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/service/oversea_pay_service"
	utility "go-oversea-pay/utility"
)

func (c *ControllerV1) Refunds(ctx context.Context, req *v1.RefundsReq) (res *v1.RefundsRes, err error) {
	currencyNumberCheck(req.Amount)
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantAccount)

	// openApiId todo mark
	_, err = oversea_pay_service.DoChannelRefund(ctx, consts.PAYMENT_BIZ_TYPE_ORDER, req, 0)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
