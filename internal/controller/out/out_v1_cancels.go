package out

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/api/out/v1"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	"go-oversea-pay/internal/logic/payment/service"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

func (c *ControllerV1) Cancels(ctx context.Context, req *v1.CancelsReq) (res *v1.CancelsRes, err error) {
	var (
		one *entity.OverseaPay
	)
	//参数有效性校验 todo mark
	merchantCheck(ctx, req.MerchantAccount)

	err = dao.OverseaPay.Ctx(ctx).Where(entity.OverseaPay{MerchantOrderNo: req.PaymentsPspReference}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil, err
	}
	utility.Assert(one != nil, "payment not found")
	err = service.DoChannelCancel(ctx, one)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
