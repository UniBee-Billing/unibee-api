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

func (c *ControllerMock) SamplePaymentNetherlands(ctx context.Context, req *mock.SamplePaymentNetherlandsReq) (res *mock.SamplePaymentNetherlandsRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, _interface.GetMerchantId(ctx))
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	outPayVo := &v12.NewPaymentReq{
		ExternalPaymentId: uuid.New().String(),
		TotalAmount:       req.Amount,
		Currency:          req.Currency,
		PaymentMethod: &v12.MethodListReq{
			Gateway: req.GatewayName,
		},
		RedirectUrl:    req.ReturnUrl,
		CountryCode:    "NL",
		Email:          "customer@email.nl",
		ExternalUserId: uuid.New().String(),
		LineItems: []*v12.OutLineItem{{
			UnitAmountExcludingTax: 22,
			Description:            uuid.New().String(),
			Quantity:               1,
		}},
	}
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig

	payments, err := NewPayment().NewPayment(ctx, outPayVo)
	if err != nil {
		return nil, err
	}
	res = &mock.SamplePaymentNetherlandsRes{
		Status:            payments.Status,
		PaymentId:         payments.PaymentId,
		MerchantPaymentId: payments.ExternalPaymentId,
		Action:            payments.Action,
	}
	return
}
