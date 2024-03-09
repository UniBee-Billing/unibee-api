package merchant

import (
	"context"
	"fmt"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/payment"
)

func (c *ControllerPayment) New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.TotalAmount > 0, "amount value is nil")
	utility.Assert(len(req.Currency) > 0, "amount currency is nil")
	currencyNumberCheck(req.TotalAmount, req.Currency)
	//utility.Assert(len(req.CountryCode) > 0, "countryCode is nil")
	utility.Assert(req.GatewayId > 0, "gatewayId is nil")
	utility.Assert(len(req.ExternalPaymentId) > 0, "ExternalPaymentId is nil")
	utility.Assert(len(req.ExternalUserId) > 0, "ExternalUserId is nil")
	utility.Assert(len(req.Email) > 0, "Email is nil")
	utility.Assert(req.Items != nil, "lineItems is nil")

	merchantInfo := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	gateway := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.MerchantId == merchantInfo.Id, "merchant gateway not match")

	var invoice *ro.InvoiceDetailSimplify
	if req.Items != nil && len(req.Items) > 0 {
		var invoiceItems []*ro.InvoiceItemDetailRo
		var totalAmountExcludingTax int64 = 0
		var totalTax int64 = 0
		for _, line := range req.Items {
			amountExcludingTax := line.UnitAmountExcludingTax * line.Quantity
			tax := int64(float64(amountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(line.TaxScale))
			invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
				Currency:               req.Currency,
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
		invoice = &ro.InvoiceDetailSimplify{
			TotalAmount:                    totalAmount,
			TotalAmountExcludingTax:        totalAmountExcludingTax,
			Currency:                       req.Currency,
			TaxAmount:                      totalTax,
			SubscriptionAmount:             totalAmount,
			SubscriptionAmountExcludingTax: totalAmountExcludingTax,
			Lines:                          invoiceItems,
		}

		utility.Assert(totalAmount == req.TotalAmount, "totalAmount not match")
	}
	user, err := auth.QueryOrCreateUser(ctx, &auth.NewReq{
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
	})
	utility.AssertError(err, "Server Error")
	utility.Assert(user != nil, "Server Error")

	createPayContext := &ro.NewPaymentInternalReq{
		CheckoutMode: true,
		Gateway:      gateway,
		Pay: &entity.Payment{
			ExternalPaymentId: req.ExternalPaymentId,
			BizType:           consts.BizTypeOneTime,
			UserId:            int64(user.Id),
			GatewayId:         gateway.Id,
			TotalAmount:       req.TotalAmount,
			Currency:          req.Currency,
			CountryCode:       req.CountryCode,
			MerchantId:        merchantInfo.Id,
			CompanyId:         merchantInfo.CompanyId,
			ReturnUrl:         req.RedirectUrl,
		},
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		Metadata:       req.Metadata,
		Invoice:        invoice,
	}

	resp, err := service.GatewayPaymentCreate(ctx, createPayContext)
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	res = &payment.NewRes{
		Status:            consts.PaymentCreated,
		PaymentId:         resp.PaymentId,
		ExternalPaymentId: req.ExternalPaymentId,
		Action:            resp.Action,
	}
	return
}
