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

func (c *ControllerMock) SamplePaymentNetherlands(ctx context.Context, req *mock.SamplePaymentNetherlandsReq) (res *mock.SamplePaymentNetherlandsRes, err error) {
	oneOpenApiConfig := query.GetOneOpenApiConfigByMerchant(ctx, req.MerchantId)
	utility.Assert(oneOpenApiConfig != nil, "openApi未设置")
	outPayVo := &v12.PaymentsReq{
		MerchantId:        req.MerchantId,
		MerchantPaymentId: uuid.New().String(),
		TotalAmount: &v12.AmountVo{
			Currency: req.Currency,
			Amount:   req.Amount,
		},
		PaymentMethod: &v12.PaymentMethodsReq{
			TokenId: "",
			Gateway: req.GatewayName,
		},
		RedirectUrl:     req.ReturnUrl,
		CountryCode:     "NL",
		TelephoneNumber: "+31689124321",
		ShopperEmail:    "customer@email.nl",
		ShopperUserId:   uuid.New().String(),
		Platform:        "WEB",
		LineItems: []*v12.OutLineItem{{
			UnitAmountExcludingTax: 22,
			Description:            uuid.New().String(),
			Quantity:               1,
		}},
		ShopperName: &v12.OutShopperName{
			FirstName: "Test",
			LastName:  "Person-nl",
		},
		//BillingAddress: &v12.OutPayAddress{
		//	City:              "Amsterdam",
		//	Country:           "NL",
		//	HouseNumberOrName: "137",
		//	PostalCode:        "1068 SR",
		//	StateOrProvince:   "33",
		//	Street:            "Osdorpplein",
		//},
		//DetailAddress: &v12.OutPayAddress{
		//	City:              "Amsterdam",
		//	Country:           "NL",
		//	HouseNumberOrName: "137",
		//	PostalCode:        "1068 SR",
		//	StateOrProvince:   "33",
		//	Street:            "Osdorpplein",
		//},
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
		Status:            payments.Status,
		PaymentId:         payments.PaymentId,
		MerchantPaymentId: payments.MerchantPaymentId,
		Action:            payments.Action,
	}
	return
}
