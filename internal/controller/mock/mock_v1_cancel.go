package mock

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	v12 "go-oversea-pay/api/out/v1"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/controller/out"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/mock/v1"
)

func (c *ControllerV1) Cancel(ctx context.Context, req *v1.CancelReq) (res *v1.CancelRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	cancelsReq := &v12.CancelsReq{
		PaymentsPspReference: req.PaymentPspReference,
		MerchantId:           req.MerchantId,
		Reference:            uuid.New().String(),
	}
	_, err = out.NewV1().Cancels(ctx, cancelsReq)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
