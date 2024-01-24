package open

import (
	"context"
	"github.com/google/uuid"
	"go-oversea-pay/api/open/mock"
	v12 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerMock) Capture(ctx context.Context, req *mock.CaptureReq) (res *mock.CaptureRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	capturesReq := &v12.CapturesReq{
		PaymentId:  req.PaymentId,
		MerchantId: req.MerchantId,
		Reference:  uuid.New().String(),
		Amount: &v12.PayAmountVo{
			Currency: req.Currency,
			Value:    req.Amount,
		},
	}
	_, err = NewPayment().Captures(ctx, capturesReq)
	if err != nil {
		return nil, err
	}
	return &mock.CaptureRes{}, nil
}
