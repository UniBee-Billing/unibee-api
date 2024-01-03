package open

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"go-oversea-pay/api/open/mock"
	v12 "go-oversea-pay/api/open/payment"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func (c *ControllerMock) Refund(ctx context.Context, req *mock.RefundReq) (res *mock.RefundRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig
	refundsReq := &v12.RefundsReq{
		PaymentsPspReference: req.PaymentPspReference,
		MerchantId:           req.MerchantId,
		Reference:            uuid.New().String(),
		Amount: &v12.PayAmountVo{
			Currency: req.Currency,
			Value:    req.Amount,
		},
	}
	_, err = NewPayment().Refunds(ctx, refundsReq)
	if err != nil {
		return nil, err
	}
	utility.SuccessJsonExit(g.RequestFromCtx(ctx), err == nil)
	return
}
