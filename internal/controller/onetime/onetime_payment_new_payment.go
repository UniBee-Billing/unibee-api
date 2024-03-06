package onetime

import (
	"context"
	"fmt"
	"unibee/api/onetime/payment"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/payment/service"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerPayment) NewPayment(ctx context.Context, req *payment.NewPaymentReq) (res *payment.NewPaymentRes, err error) {
	utility.Assert(req != nil, "request req is nil")
	utility.Assert(req.TotalAmount != nil, "amount is nil")
	utility.Assert(req.TotalAmount.Amount > 0, "amount value is nil")
	utility.Assert(len(req.TotalAmount.Currency) > 0, "amount currency is nil")
	currencyNumberCheck(req.TotalAmount)
	utility.Assert(len(req.CountryCode) > 0, "countryCode is nil")
	utility.Assert(req.PaymentMethod != nil, "payment method is nil")
	utility.Assert(len(req.PaymentMethod.Gateway) > 0, "payment method type is nil")
	utility.Assert(len(req.ExternalPaymentId) > 0, "ExternalPaymentId is nil")
	utility.Assert(len(req.ExternalUserId) > 0, "shopperReference type is nil")
	utility.Assert(len(req.Email) > 0, "shopperEmail is nil")
	utility.Assert(req.LineItems != nil, "lineItems is nil")

	_, merchantInfo := merchantCheck(ctx, _interface.GetMerchantId(ctx))
	gateway := query.GetGatewayByGatewayName(ctx, merchantInfo.Id, req.PaymentMethod.Gateway)
	utility.Assert(gateway != nil, "type not found:"+req.PaymentMethod.Gateway)

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

	createPayContext := &ro.NewPaymentInternalReq{
		Gateway: gateway,
		Pay: &entity.Payment{
			ExternalPaymentId: req.ExternalPaymentId,
			BizType:           consts.BizTypeOneTime,
			GatewayId:         gateway.Id,
			TotalAmount:       req.TotalAmount.Amount,
			Currency:          req.TotalAmount.Currency,
			CountryCode:       req.CountryCode,
			MerchantId:        merchantInfo.Id,
			CompanyId:         merchantInfo.CompanyId,
			ReturnUrl:         req.RedirectUrl,
		},
		ExternalUserId: req.ExternalUserId,
		Email:          req.Email,
		MetaData:       req.Metadata,
		Invoice:        invoice,
	}

	resp, err := service.GatewayPaymentCreate(ctx, createPayContext)
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	res = &payment.NewPaymentRes{
		Status:            "Pending",
		PaymentId:         resp.PaymentId,
		ExternalPaymentId: req.ExternalPaymentId,
		Action:            resp.Action,
	}
	return
}
