package invoice_compute

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func GetInvoiceLink(invoiceId string) string {
	return fmt.Sprintf("%s/in/%s", consts.GetConfigInstance().Server.GetServerPath(), invoiceId)
}

func ConvertInvoiceToRo(ctx context.Context, invoice *entity.Invoice) *bean.InvoiceDetail {
	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(invoice.Lines, &lines)
	for _, line := range lines {
		line.Currency = invoice.Currency
		line.TaxScale = invoice.TaxScale
	}
	if err != nil {
		fmt.Printf("ConvertInvoiceLines err:%s", err)
	}
	return &bean.InvoiceDetail{
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
		GatewayId:                      invoice.GatewayId,
		Status:                         invoice.Status,
		SendStatus:                     invoice.SendStatus,
		SendEmail:                      invoice.SendEmail,
		SendPdf:                        invoice.SendPdf,
		UserId:                         invoice.UserId,
		Data:                           invoice.Data,
		GmtModify:                      invoice.GmtModify,
		IsDeleted:                      invoice.IsDeleted,
		Link:                           invoice.Link,
		GatewayStatus:                  invoice.GatewayStatus,
		GatewayInvoiceId:               invoice.GatewayInvoiceId,
		GatewayInvoicePdf:              invoice.GatewayInvoicePdf,
		TaxScale:                       invoice.TaxScale,
		SendNote:                       invoice.SendNote,
		SendTerms:                      invoice.SendTerms,
		DiscountAmount:                 0,
		TotalAmountExcludingTax:        invoice.TotalAmountExcludingTax,
		SubscriptionAmountExcludingTax: invoice.SubscriptionAmountExcludingTax,
		PeriodStart:                    invoice.PeriodStart,
		PeriodEnd:                      invoice.PeriodEnd,
		Gateway:                        bean.SimplifyGateway(query.GetGatewayById(ctx, invoice.GatewayId)),
		Merchant:                       bean.SimplifyMerchant(query.GetMerchantById(ctx, invoice.MerchantId)),
		UserAccount:                    bean.SimplifyUserAccount(query.GetUserAccountById(ctx, uint64(invoice.UserId))),
		Subscription:                   bean.SimplifySubscription(query.GetSubscriptionBySubscriptionId(ctx, invoice.SubscriptionId)),
		Payment:                        bean.SimplifyPayment(query.GetPaymentByPaymentId(ctx, invoice.PaymentId)),
		Refund:                         bean.SimplifyRefund(query.GetRefundByRefundId(ctx, invoice.RefundId)),
	}
}

type CalculateInvoiceReq struct {
	Currency      string `json:"currency"`
	PlanId        uint64 `json:"planId"`
	Quantity      int64  `json:"quantity"`
	AddonJsonData string `json:"addonJsonData"`
	TaxScale      int64  `json:"taxScale"`
	PeriodStart   int64  `json:"periodStart"`
	PeriodEnd     int64  `json:"periodEnd"`
	InvoiceName   string `json:"invoiceName"`
}

func ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx context.Context, req *CalculateInvoiceReq) *bean.InvoiceSimplify {
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, fmt.Sprintf("plan not found:%d", req.PlanId))
	addons := addon2.GetSubscriptionAddonsByAddonJson(ctx, req.AddonJsonData)
	var totalAmountExcludingTax = plan.Amount * req.Quantity
	for _, addon := range addons {
		totalAmountExcludingTax = totalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var period = ""
	if req.PeriodStart > 0 && req.PeriodEnd > req.PeriodStart {
		period = fmt.Sprintf("(%s-%s)", gtime.NewFromTimeStamp(req.PeriodStart).Layout("2006-01-02"), gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02"))
	}

	var invoiceItems []*bean.InvoiceItemSimplify
	invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
		Currency:               req.Currency,
		Amount:                 req.Quantity*plan.Amount + int64(float64(req.Quantity*plan.Amount)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
		AmountExcludingTax:     req.Quantity * plan.Amount,
		Tax:                    int64(float64(req.Quantity*plan.Amount) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
		TaxScale:               req.TaxScale,
		UnitAmountExcludingTax: plan.Amount,
		Description:            fmt.Sprintf("%d * %s %s", req.Quantity, plan.PlanName, period),
		Quantity:               req.Quantity,
	})
	for _, addon := range addons {
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			Amount:                 addon.Quantity*addon.AddonPlan.Amount + int64(float64(addon.Quantity*addon.AddonPlan.Amount)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			Tax:                    int64(float64(addon.Quantity*addon.AddonPlan.Amount) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			TaxScale:               req.TaxScale,
			AmountExcludingTax:     addon.Quantity * addon.AddonPlan.Amount,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Description:            fmt.Sprintf("%d * %s %s", addon.Quantity, addon.AddonPlan.PlanName, period),
			Quantity:               addon.Quantity,
		})
	}
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale))
	return &bean.InvoiceSimplify{
		InvoiceName:                    req.InvoiceName,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       req.Currency,
		TaxAmount:                      taxAmount,
		TaxScale:                       req.TaxScale,
		SubscriptionAmount:             totalAmountExcludingTax + taxAmount, // 在没有 discount 之前，保持于 Total 一致
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,             // 在没有 discount 之前，保持于 Total 一致
		Lines:                          invoiceItems,
		PeriodStart:                    req.PeriodStart,
		PeriodEnd:                      req.PeriodEnd,
	}
}

type ProrationPlanParam struct {
	PlanId   uint64
	Quantity int64
}

type CalculateProrationInvoiceReq struct {
	Currency          string                `json:"currency"`
	TaxScale          int64                 `json:"taxScale"`
	ProrationDate     int64                 `json:"prorationStart"`
	PeriodStart       int64                 `json:"periodStart"`
	PeriodEnd         int64                 `json:"periodEnd"`
	OldProrationPlans []*ProrationPlanParam `json:"oldPlans"`
	NewProrationPlans []*ProrationPlanParam `json:"newPlans"`
	InvoiceName       string                `json:"invoiceName"`
}

func ComputeSubscriptionProrationInvoiceDetailSimplify(ctx context.Context, req *CalculateProrationInvoiceReq) *bean.InvoiceSimplify {
	if req.OldProrationPlans == nil {
		req.OldProrationPlans = make([]*ProrationPlanParam, 0)
	}
	if req.NewProrationPlans == nil {
		req.NewProrationPlans = make([]*ProrationPlanParam, 0)
	}
	newMap := make(map[uint64]*ProrationPlanParam)
	//oldMap := make(map[int64]*ProrationPlanParam)
	//for _, planSub := range req.OldProrationPlans {
	//	oldMap[planSub.PlanId] = planSub
	//}
	for _, planSub := range req.NewProrationPlans {
		newMap[planSub.PlanId] = planSub
	}

	utility.Assert(req.ProrationDate > 0, "Invalid ProrationDate")
	utility.Assert(req.PeriodStart <= req.ProrationDate && req.ProrationDate <= req.PeriodEnd, "System Error, Subscription Need Update")

	timeScale := int64((float64(req.PeriodEnd-req.ProrationDate) / float64(req.PeriodEnd-req.PeriodStart)) * 10000)
	var invoiceItems []*bean.InvoiceItemSimplify
	var totalAmountExcludingTax int64
	for _, oldPlanSub := range req.OldProrationPlans {
		plan := query.GetPlanById(ctx, oldPlanSub.PlanId)
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(oldPlanSub.PlanId, 10))
		unitAmountExcludingTax := int64(float64(plan.Amount) * utility.ConvertTaxScaleToInternalFloat(timeScale))
		if newPlanSub, ok := newMap[oldPlanSub.PlanId]; ok {
			//new plan contain old
			quantityDiff := newPlanSub.Quantity - oldPlanSub.Quantity
			if quantityDiff > 0 {
				// quantity increase
				invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
					Currency:               req.Currency,
					Amount:                 quantityDiff*unitAmountExcludingTax + int64(float64(quantityDiff*unitAmountExcludingTax)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
					AmountExcludingTax:     quantityDiff * unitAmountExcludingTax,
					Tax:                    int64(float64(quantityDiff*unitAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
					TaxScale:               req.TaxScale,
					UnitAmountExcludingTax: unitAmountExcludingTax,
					Description:            fmt.Sprintf("Remaining Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.ProrationDate).Layout("2006-01-02")),
					Quantity:               quantityDiff,
				})
				totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
			} else if quantityDiff < 0 {
				// quantity decrease
				quantityDiff = -quantityDiff
				unitAmountExcludingTax = -unitAmountExcludingTax
				invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
					Currency:               req.Currency,
					Amount:                 quantityDiff*unitAmountExcludingTax + int64(float64(quantityDiff*unitAmountExcludingTax)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
					AmountExcludingTax:     quantityDiff * unitAmountExcludingTax,
					Tax:                    int64(float64(quantityDiff*unitAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
					TaxScale:               req.TaxScale,
					UnitAmountExcludingTax: unitAmountExcludingTax,
					Description:            fmt.Sprintf("Unused Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02")),
					Quantity:               quantityDiff,
				})
				totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
			}
			delete(newMap, newPlanSub.PlanId)
		} else {
			//old removed
			quantityDiff := oldPlanSub.Quantity
			unitAmountExcludingTax = -unitAmountExcludingTax
			invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
				Currency:               req.Currency,
				Amount:                 quantityDiff*unitAmountExcludingTax + int64(float64(quantityDiff*unitAmountExcludingTax)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
				AmountExcludingTax:     quantityDiff * unitAmountExcludingTax,
				Tax:                    int64(float64(quantityDiff*unitAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
				TaxScale:               req.TaxScale,
				UnitAmountExcludingTax: unitAmountExcludingTax,
				Description:            fmt.Sprintf("Unused Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02")),
				Quantity:               quantityDiff,
			})
			totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
		}
	}
	for _, newPlanSub := range newMap {
		plan := query.GetPlanById(ctx, newPlanSub.PlanId)
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(newPlanSub.PlanId, 10))
		unitAmountExcludingTax := int64(float64(plan.Amount) * utility.ConvertTaxScaleToInternalFloat(timeScale))
		quantityDiff := newPlanSub.Quantity
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			Amount:                 quantityDiff*unitAmountExcludingTax + int64(float64(quantityDiff*unitAmountExcludingTax)*utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			AmountExcludingTax:     quantityDiff * unitAmountExcludingTax,
			Tax:                    int64(float64(quantityDiff*unitAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale)),
			TaxScale:               req.TaxScale,
			UnitAmountExcludingTax: unitAmountExcludingTax,
			Description:            fmt.Sprintf("Remaining Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.ProrationDate).Layout("2006-01-02")),
			Quantity:               quantityDiff,
		})
		totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
	}

	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxScaleToInternalFloat(req.TaxScale))
	return &bean.InvoiceSimplify{
		InvoiceName:                    req.InvoiceName,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		Currency:                       req.Currency,
		TaxAmount:                      taxAmount,
		TaxScale:                       req.TaxScale,
		SubscriptionAmount:             totalAmountExcludingTax + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax,
		Lines:                          invoiceItems,
		ProrationDate:                  req.ProrationDate,
		ProrationScale:                 timeScale,
		PeriodStart:                    req.ProrationDate,
		PeriodEnd:                      req.PeriodEnd,
	}
}
