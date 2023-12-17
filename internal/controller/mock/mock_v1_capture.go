package mock

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	v12 "go-oversea-pay/api/out/v1"
	"go-oversea-pay/api/out/vo"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/controller/out"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/mock/v1"
)

func (c *ControllerV1) Capture(ctx context.Context, req *v1.CaptureReq) (res *v1.CaptureRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	capturesReq := &v12.CapturesReq{
		PaymentsPspReference: req.PaymentPspReference,
		MerchantId:           req.MerchantId,
		Reference:            uuid.New().String(),
		Amount: &vo.PayAmountVo{
			Currency: req.Currency,
			Value:    req.Amount,
		},
	}
	_, err = out.NewV1().Captures(ctx, capturesReq)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
