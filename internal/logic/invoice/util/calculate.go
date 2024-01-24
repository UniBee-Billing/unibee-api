package util

import (
	"context"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

type CalculateInvoiceReq struct {
	Currency      string `json:"currency"`
	PlanId        int64  `json:"planId"`
	Quantity      int64  `json:"quantity"`
	AddonJsonData string `json:"addonJsonData"`
	TaxScale      int64  `json:"taxScale"`
}

func CalculateInternalInvoiceRo(ctx context.Context, req *CalculateInvoiceReq) *ro.ChannelDetailInvoiceRo {
	plan := query.GetPlanById(ctx, req.PlanId)
	addons := query.GetSubscriptionAddonsByAddonJson(ctx, req.AddonJsonData)
	var totalAmountExcludingTax = plan.Amount * req.Quantity
	for _, addon := range addons {
		totalAmountExcludingTax = totalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var invoiceItems []*ro.ChannelDetailInvoiceItem
	invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
		Currency:               req.Currency,
		Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxPercentageToInternalFloat(req.TaxScale)),
		AmountExcludingTax:     req.Quantity * plan.Amount,
		Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(req.TaxScale)),
		UnitAmountExcludingTax: plan.Amount,
		Description:            plan.PlanName,
		Quantity:               req.Quantity,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &ro.ChannelDetailInvoiceItem{
			Currency:               req.Currency,
			Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxPercentageToInternalFloat(req.TaxScale)),
			Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxPercentageToInternalFloat(req.TaxScale)),
			AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Description:            addon.AddonPlan.PlanName,
			Quantity:               addon.Quantity,
		})
	}
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxScale))
	return &ro.ChannelDetailInvoiceRo{
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       req.Currency,
		TaxAmount:                      taxAmount,
		SubscriptionAmount:             totalAmountExcludingTax + taxAmount, // 在没有 discount 之前，保持于 Total 一致
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,             // 在没有 discount 之前，保持于 Total 一致
		Lines:                          invoiceItems,
	}
}
