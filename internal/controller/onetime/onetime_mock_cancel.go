package onetime

import (
	"context"
	"github.com/google/uuid"
	"unibee-api/api/onetime/mock"
	v12 "unibee-api/api/onetime/payment"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerMock) Cancel(ctx context.Context, req *mock.CancelReq) (res *mock.CancelRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	cancelsReq := &v12.CancelReq{
		PaymentId:        req.PaymentId,
		MerchantId:       req.MerchantId,
		MerchantCancelId: uuid.New().String(),
	}
	_, err = NewPayment().Cancel(ctx, cancelsReq)
	if err != nil {
		return nil, err
	}
	return &mock.CancelRes{}, nil
}
