package util

import (
	"context"
	"fmt"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func ConvertInvoiceToRo(invoice *entity.Invoice) *ro.InvoiceDetailRo {
	var lines []*ro.InvoiceItemDetailRo
	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
	for _, line := range lines {
		line.Currency = invoice.Currency
		line.TaxScale = invoice.TaxScale
	}
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	return &ro.InvoiceDetailRo{
		Id:                             invoice.Id,
		MerchantId:                     invoice.MerchantId,
		SubscriptionId:                 invoice.SubscriptionId,
		InvoiceId:                      invoice.InvoiceId,
		InvoiceName:                    invoice.InvoiceName,
		GmtCreate:                      invoice.GmtCreate,
		TotalAmount:                    invoice.TotalAmount,
		TaxAmount:                      invoice.TaxAmount,
		SubscriptionAmount:             invoice.SubscriptionAmount,
		Currency:                       invoice.Currency,
		Lines:                          lines,
		ChannelId:                      invoice.ChannelId,
		Status:                         invoice.Status,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      invoice.SendEmail,
		SendPdf:                        invoice.SendPdf,
		UserId:                         invoice.UserId,
		Data:                           invoice.Data,
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           invoice.Link,
		ChannelStatus:                  invoice.ChannelStatus,
		ChannelInvoiceId:               invoice.ChannelInvoiceId,
		ChannelInvoicePdf:              invoice.ChannelInvoicePdf,
		TaxScale:                       invoice.TaxScale,
		SendNote:                       invoice.SendNote,
		SendTerms:                      invoice.SendTerms,
		DiscountAmount:                 0,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
	}
}

type CalculateInvoiceReq struct {
	Currency      string `json:"currency"`
	PlanId        int64  `json:"planId"`
	Quantity      int64  `json:"quantity"`
	AddonJsonData string `json:"addonJsonData"`
	TaxScale      int64  `json:"taxScale"`
}

func ConvertInvoiceDetailSimplifyFromPlanAddons(ctx context.Context, req *CalculateInvoiceReq) *ro.InvoiceDetailSimplify {
	plan := query.GetPlanById(ctx, req.PlanId)
	addons := query.GetSubscriptionAddonsByAddonJson(ctx, req.AddonJsonData)
	var totalAmountExcludingTax = plan.Amount * req.Quantity
	for _, addon := range addons {
		totalAmountExcludingTax = totalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var invoiceItems []*ro.InvoiceItemDetailRo
	invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
		Currency:               req.Currency,
		Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
		AmountExcludingTax:     req.Quantity * plan.Amount,
		Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
		UnitAmountExcludingTax: plan.Amount,
		Description:            plan.PlanName,
		Quantity:               req.Quantity,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               req.Currency,
			Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Description:            addon.AddonPlan.PlanName,
			Quantity:               addon.Quantity,
		})
	}
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale))
	return &ro.InvoiceDetailSimplify{
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       req.Currency,
		TaxAmount:                      taxAmount,
		SubscriptionAmount:             totalAmountExcludingTax + taxAmount, // 在没有 discount 之前，保持于 Total 一致
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,             // 在没有 discount 之前，保持于 Total 一致
		Lines:                          invoiceItems,
	}
}
