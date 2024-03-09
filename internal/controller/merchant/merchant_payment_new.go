package merchant

import (
	"context"
	"fmt"
	"strings"
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
	req.Currency = strings.ToUpper(req.Currency)

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
			utility.Assert(line.Amount > 0, "Item Amount invalid, should > 0")
			utility.Assert(len(line.Description) > 0, "Item Description invalid")
			invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
				Currency:               line.Currency,
				TaxScale:               line.TaxScale,
				Tax:                    line.Tax,
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: line.UnitAmountExcludingTax,
				Description:            line.Description,
				Quantity:               line.Quantity,
			})
			totalTax = totalTax + line.Tax
			totalAmountExcludingTax = totalAmountExcludingTax + (line.Amount - line.Tax)
		}
		var totalAmount = totalTax + totalAmountExcludingTax
		utility.Assert(totalAmount == req.TotalAmount, "sum(items.amount) should match totalAmount")
		invoice = &ro.InvoiceDetailSimplify{
			TotalAmount:             req.TotalAmount,
			Currency:                req.Currency,
			TotalAmountExcludingTax: totalAmountExcludingTax,
			TaxAmount:               totalTax,
			Lines:                   invoiceItems,
		}
	} else {
		invoice = &ro.InvoiceDetailSimplify{
			TotalAmount:             req.TotalAmount,
			TotalAmountExcludingTax: req.TotalAmount,
			Currency:                req.Currency,
			TaxAmount:               0,
			Lines: []*ro.InvoiceItemDetailRo{{
				Currency:               req.Currency,
				Amount:                 req.TotalAmount,
				Tax:                    0,
				AmountExcludingTax:     req.TotalAmount,
				TaxScale:               0,
				UnitAmountExcludingTax: req.TotalAmount,
				Description:            merchantInfo.Name,
				Quantity:               1,
			}},
		}
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
		Link:              resp.Link,
		Action:            resp.Action,
	}
	return
}
