package onetime

import (
	"context"
	"github.com/google/uuid"
	"unibee/api/onetime/mock"
	v12 "unibee/api/onetime/payment"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerMock) Capture(ctx context.Context, req *mock.CaptureReq) (res *mock.CaptureRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	capturesReq := &v12.CaptureReq{
		PaymentId:         req.PaymentId,
		MerchantCaptureId: uuid.New().String(),
		Amount: &v12.AmountVo{
			Currency: req.Currency,
			Amount:   req.Amount,
		},
	}
	_, err = NewPayment().Capture(ctx, capturesReq)
	if err != nil {
		return nil, err
	}
	return &mock.CaptureRes{}, nil
}
