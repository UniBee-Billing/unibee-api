package invoice_compute

import (
	"context"
	"unibee/api/bean"
	"unibee/internal/logic/discount"
	"unibee/utility"
)

type InvoiceSimplifyInternalReq struct {
	Id             uint64                            `json:"id"`
	MerchantId     uint64                            `json:"merchantId"`
	InvoiceId      string                            `json:"invoiceId"`
	InvoiceName    string                            `json:"invoiceName"`
	DiscountCode   string                            `json:"discountCode"`
	Currency       string                            `json:"currency"`
	TaxPercentage  int64                             `json:"taxPercentage"`
	Lines          []*InvoiceItemSimplifyInternalReq `json:"lines"`
	PeriodEnd      int64                             `json:"periodEnd"`
	PeriodStart    int64                             `json:"periodStart"`
	ProrationDate  int64                             `json:"prorationDate"`
	ProrationScale int64                             `json:"prorationScale"`
	FinishTime     int64                             `json:"finishTime"`
	SendStatus     int                               `json:"sendStatus"`
	DayUtilDue     int64                             `json:"dayUtilDue"`
	TimeNow        int64                             `json:"timeNow"`
}

type InvoiceItemSimplifyInternalReq struct {
	UnitAmountExcludingTax int64              `json:"unitAmountExcludingTax"`
	Quantity               int64              `json:"quantity"`
	Description            string             `json:"description"`
	Plan                   *bean.PlanSimplify `json:"plan"`
}

func MakeInvoiceSimplify(ctx context.Context, req *InvoiceSimplifyInternalReq) *bean.InvoiceSimplify {
	utility.Assert(req.Lines != nil, "MakeInvoiceSimplify error, line is null")
	utility.Assert(req.MerchantId > 0, "MakeInvoiceSimplify error, merchantId is null")
	var invoiceItems = make([]*bean.InvoiceItemSimplify, 0)
	var totalAmountExcludingTax int64 = 0
	for _, item := range req.Lines {
		var amountExcludingTax = item.Quantity * item.UnitAmountExcludingTax
		var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			OriginAmount:           amountExcludingTax + taxAmount,
			Amount:                 amountExcludingTax + taxAmount,
			Tax:                    taxAmount,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: item.UnitAmountExcludingTax,
			Quantity:               item.Quantity,
			Description:            item.Description,
			Plan:                   item.Plan,
		})
		totalAmountExcludingTax = totalAmountExcludingTax + amountExcludingTax
	}

	discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, req.MerchantId, totalAmountExcludingTax, req.Currency, req.DiscountCode, req.TimeNow), totalAmountExcludingTax)
	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
	ProrationDiscountToItem(discountAmount, taxAmount, invoiceItems)

	return &bean.InvoiceSimplify{
		InvoiceName:                    req.InvoiceName,
		OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		DiscountAmount:                 discountAmount,
		DiscountCode:                   req.DiscountCode,
		TaxAmount:                      taxAmount,
		Currency:                       req.Currency,
		TaxPercentage:                  req.TaxPercentage,
		SubscriptionAmount:             totalAmountExcludingTax + discountAmount + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount,
		Lines:                          invoiceItems,
		PeriodStart:                    req.PeriodStart,
		PeriodEnd:                      req.PeriodEnd,
		FinishTime:                     req.FinishTime,
		SendStatus:                     req.SendStatus,
		DayUtilDue:                     req.DayUtilDue,
	}
}
