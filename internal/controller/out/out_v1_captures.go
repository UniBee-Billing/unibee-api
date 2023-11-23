package out

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/service/oversea_pay_service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/out/v1"
)

func (c *ControllerV1) Captures(ctx context.Context, req *v1.CapturesReq) (res *v1.CapturesRes, err error) {
	var (
		one *entity.OverseaPay
	)
	err = dao.OverseaPay.Ctx(ctx).Where(entity.OverseaPay{MerchantOrderNo: req.PaymentsPspReference}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil, err
	}
	utility.Assert(one != nil, "payment not found")
	utility.Assert(one.Currency == req.Amount.Currency, "Currency not match the payment")
	one.BuyerPayFee = req.Amount.Value
	result, err := oversea_pay_service.DoChannelCapture(ctx, one)
	utility.Assert(err == nil, err.Error())
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), result)
	return
}
