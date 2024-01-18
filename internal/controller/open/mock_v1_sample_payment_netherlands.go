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

func (c *ControllerMock) SamplePaymentNetherlands(ctx context.Context, req *mock.SamplePaymentNetherlandsReq) (res *mock.SamplePaymentNetherlandsRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	outPayVo := &v12.PaymentsReq{
		MerchantId: req.MerchantId,
		Reference:  uuid.New().String(),
		Amount: &v12.PayAmountVo{
			Currency: req.Currency,
			Value:    req.Amount,
		},
		PaymentMethod: &v12.PaymentMethodsReq{
			TokenId: "",
			Channel: req.Channel,
		},
		PaymentBrandAddition: nil,
		StorePaymentMethod:   false,
		ReturnUrl:            req.ReturnUrl,
		CountryCode:          "NL",
		TelephoneNumber:      "+31689124321",
		ShopperEmail:         "customer@email.nl",
		ShopperReference:     uuid.New().String(),
		Channel:              "WEB",
		LineItems: []*v12.OutLineItem{{
			AmountExcludingTax: 22,
			AmountIncludingTax: 11,
			Description:        uuid.New().String(),
			Id:                 uuid.New().String(),
			Quantity:           1,
			TaxAmount:          11,
		}},
		ShopperName: &v12.OutShopperName{
			FirstName: "Test",
			LastName:  "Person-nl",
		},
		BillingAddress: &v12.OutPayAddress{
			City:              "Amsterdam",
			Country:           "NL",
			HouseNumberOrName: "137",
			PostalCode:        "1068 SR",
			StateOrProvince:   "33",
			Street:            "Osdorpplein",
		},
		DetailAddress: &v12.OutPayAddress{
			City:              "Amsterdam",
			Country:           "NL",
			HouseNumberOrName: "137",
			PostalCode:        "1068 SR",
			StateOrProvince:   "33",
			Street:            "Osdorpplein",
		},
		Capture:           false,
		CaptureDelayHours: 0,
		DateOfBrith:       "1970-10-07",
	}
	_interface.BizCtx().Get(ctx).Data[consts.ApiKey] = oneOpenApiConfig.ApiKey
	_interface.BizCtx().Get(ctx).OpenApiConfig = oneOpenApiConfig

	payments, err := NewPayment().Payments(ctx, outPayVo)
	if err != nil {
		return nil, err
	}
	res = &mock.SamplePaymentNetherlandsRes{
		Status:    payments.Status,
		PaymentId: payments.PaymentId,
		Reference: payments.Reference,
		Action:    payments.Action,
	}
	return
}
