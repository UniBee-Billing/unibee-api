package open

import (
	"context"
	"github.com/google/uuid"
	"unibee-api/api/open/mock"
	v12 "unibee-api/api/open/payment"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerMock) Refund(ctx context.Context, req *mock.RefundReq) (res *mock.RefundRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	refundsReq := &v12.RefundsReq{
		PaymentId:        req.PaymentId,
		MerchantId:       req.MerchantId,
		MerchantRefundId: uuid.New().String(),
		Amount: &v12.AmountVo{
			Currency: req.Currency,
			Amount:   req.Amount,
		},
	}
	_, err = NewPayment().Refunds(ctx, refundsReq)
	if err != nil {
		return nil, err
	}
	return &mock.RefundRes{}, nil
}
