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

func (c *ControllerMock) Cancel(ctx context.Context, req *mock.CancelReq) (res *mock.CancelRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	cancelsReq := &v12.CancelsReq{
		PaymentId:  req.PaymentId,
		MerchantId: req.MerchantId,
		Reference:  uuid.New().String(),
	}
	_, err = NewPayment().Cancels(ctx, cancelsReq)
	if err != nil {
		return nil, err
	}
	return &mock.CancelRes{}, nil
}
