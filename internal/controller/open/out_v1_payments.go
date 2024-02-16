package open

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee-api/api/open/payment"
	"unibee-api/internal/consts"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/payment/service"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerPayment) Payments(ctx context.Context, req *payment.PaymentsReq) (res *payment.PaymentsRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.TotalAmount != nil, "amount is nil")
	utility.Assert(req.TotalAmount.Amount > 0, "amount value is nil")
	utility.Assert(len(req.TotalAmount.Currency) > 0, "amount currency is nil")
	//类似日元的小数尾数必须为 0 检查
	currencyNumberCheck(req.TotalAmount)
	utility.Assert(len(req.CountryCode) > 0, "countryCode is nil")
	utility.Assert(req.MerchantId > 0, "merchantId is nil")
	utility.Assert(req.PaymentMethod != nil, "payment method is nil")
	utility.Assert(len(req.PaymentMethod.Gateway) > 0, "payment method type is nil")
	utility.Assert(len(req.MerchantPaymentId) > 0, "MerchantPaymentId is nil")
	utility.Assert(len(req.ShopperUserId) > 0, "shopperReference type is nil")
	utility.Assert(len(req.ShopperEmail) > 0, "shopperEmail is nil")
	utility.Assert(req.LineItems != nil, "lineItems is nil")

	openApiConfig, merchantInfo := merchantCheck(ctx, req.MerchantId)
	gateway := query.GetGatewayByGatewayName(ctx, req.PaymentMethod.Gateway)
	utility.Assert(gateway != nil, "type not found:"+req.PaymentMethod.Gateway)
	//支付方式绑定校验 todo mark

	var invoiceItems []*ro.InvoiceItemDetailRo
	var totalAmountExcludingTax int64 = 0
	var totalTax int64 = 0
	for _, line := range req.LineItems {
		amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
		tax := int64(float64(amountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(line.TaxScale))
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               req.TotalAmount.Currency,
			TaxScale:               line.TaxScale,
			Tax:                    tax,
			Amount:                 amountExcludingTax + tax,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: line.UnitAmountExcludingTax,
			Description:            line.Description,
			Quantity:               line.Quantity,
		})
		totalTax = totalTax + tax
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}
	var totalAmount = totalTax + totalAmountExcludingTax
	var invoice = &ro.InvoiceDetailSimplify{
		TotalAmount:                    totalAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       req.TotalAmount.Currency,
		TaxAmount:                      totalTax,
		SubscriptionAmount:             totalAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		Lines:                          invoiceItems,
	}

	utility.Assert(totalAmount == req.TotalAmount.Amount, "totalAmount not match")

	createPayContext := &ro.CreatePayContext{
		OpenApiId: int64(openApiConfig.Id),
		Gateway:   gateway,
		Pay: &entity.Payment{
			BizId:             req.MerchantPaymentId,
			BizType:           consts.BIZ_TYPE_ONE_TIME,
			GatewayId:         int64(gateway.Id),
			TotalAmount:       req.TotalAmount.Amount,
			Currency:          req.TotalAmount.Currency,
			CountryCode:       req.CountryCode,
			MerchantId:        merchantInfo.Id,
			CompanyId:         merchantInfo.CompanyId,
			ReturnUrl:         req.RedirectUrl,
			CaptureDelayHours: req.CaptureDelayHours,
		},
		Platform:      req.Platform,
		DeviceType:    req.DeviceType,
		ShopperUserId: req.ShopperUserId,
		ShopperEmail:  req.ShopperEmail,
		ShopperLocale: req.ShopperLocale,
		Mobile:        req.TelephoneNumber,
		MediaData:     req.Metadata,
		Invoice:       invoice,
		//BillingDetails:           req.BillingAddress,
		//ShippingDetails:          req.DetailAddress,
		ShopperName:              req.ShopperName,
		ShopperInteraction:       req.ShopperInteraction,
		RecurringProcessingModel: req.RecurringProcessingToken,
		TokenId:                  req.PaymentMethod.TokenId,
		MerchantOrderReference:   req.MerchantOrderReference,
		DateOfBirth:              gtime.ParseTimeFromContent(req.DateOfBrith, "YYYY-MM-DD"),
		PayMethod:                1, //automatic
		DaysUtilDue:              5, //one day expire
	}

	resp, err := service.GatewayPaymentCreate(ctx, createPayContext)
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	res = &payment.PaymentsRes{
		Status:            "Pending",
		PaymentId:         resp.PaymentId,
		MerchantPaymentId: req.MerchantPaymentId,
		Action:            resp.Action,
	}
	return
}
