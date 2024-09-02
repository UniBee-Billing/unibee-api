package invoice_compute

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"math"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	"unibee/internal/logic/discount"
	subscription2 "unibee/internal/logic/subscription"
	addon2 "unibee/internal/logic/subscription/addon"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type CalculateInvoiceReq struct {
	Currency           string                 `json:"currency"`
	DiscountCode       string                 `json:"discountCode"`
	TimeNow            int64                  `json:"TimeNow"`
	PlanId             uint64                 `json:"planId"`
	Quantity           int64                  `json:"quantity"`
	AddonJsonData      string                 `json:"addonJsonData"`
	CountryCode        string                 `json:"CountryCode"`
	VatNumber          string                 `json:"vatNumber"`
	TaxPercentage      int64                  `json:"taxPercentage"`
	PeriodStart        int64                  `json:"periodStart"`
	PeriodEnd          int64                  `json:"periodEnd"`
	FinishTime         int64                  `json:"finishTime"`
	InvoiceName        string                 `json:"invoiceName"`
	ProductData        *bean.PlanProductParam `json:"productData"  dc:"ProductData"  `
	BillingCycleAnchor int64                  `json:"billingCycleAnchor"             description:"billing_cycle_anchor"` // billing_cycle_anchor
	CreateFrom         string                 `json:"createFrom"                     description:"create from"`          // create from
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func VerifyInvoiceSimplify(one *bean.Invoice) {
	var totalAmount = one.TotalAmount
	var totalOriginAmount = one.OriginAmount
	var totalTax = one.TaxAmount
	var totalDiscountAmount = one.DiscountAmount
	for _, item := range one.Lines {
		totalAmount = totalAmount - item.Amount
		totalOriginAmount = totalOriginAmount - item.OriginAmount
		totalTax = totalTax - item.Tax
		totalDiscountAmount = totalDiscountAmount - item.DiscountAmount
		utility.Assert(one.TaxPercentage == item.TaxPercentage, "taxPercentage is not match")
		utility.Assert(item.AmountExcludingTax == item.UnitAmountExcludingTax*item.Quantity-item.DiscountAmount, "item AmountExcludingTax not match unit*quantity-discount")
		utility.Assert(one.Currency == item.Currency, "currency not match")
	}
	utility.Assert(totalAmount == 0, "totalAmount is not equal to lines")
	utility.Assert(totalOriginAmount == 0, "totalOriginAmount is not equal to lines")
	utility.Assert(totalTax == 0, "totalTax is not equal to lines")
	utility.Assert(totalDiscountAmount == 0, "totalDiscountAmount is not equal to lines")
	if one.Status >= consts.InvoiceStatusProcessing {
		utility.Assert(one.FinishTime != 0, "process invoice has no finishTime")
	}

}
func VerifyInvoice(one *entity.Invoice) {
	var lines []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(one.Lines, &lines)
	utility.AssertError(err, "VerifyInvoice")
	var totalAmount = one.TotalAmount
	var totalOriginAmount = one.TotalAmount + one.TaxAmount
	var totalTax = one.TaxAmount
	var totalDiscountAmount = one.DiscountAmount
	for _, item := range lines {
		totalAmount = totalAmount - item.Amount
		totalOriginAmount = totalOriginAmount - item.OriginAmount
		totalTax = totalTax - item.Tax
		totalDiscountAmount = totalDiscountAmount - item.DiscountAmount
		utility.Assert(one.TaxPercentage == item.TaxPercentage, "taxPercentage is not match")
		utility.Assert(item.AmountExcludingTax == item.UnitAmountExcludingTax*item.Quantity-item.DiscountAmount, "item AmountExcludingTax not match unit*quantity-discount")
		utility.Assert(one.Currency == item.Currency, "currency not match")
	}
	utility.Assert(totalAmount == 0, "totalAmount is not equal to lines")
	utility.Assert(totalOriginAmount == 0, "totalOriginAmount is not equal to lines")
	utility.Assert(totalTax == 0, "totalTax is not equal to lines")
	utility.Assert(totalDiscountAmount == 0, "totalDiscountAmount is not equal to lines")
	if one.Status >= consts.InvoiceStatusProcessing {
		utility.Assert(one.FinishTime != 0, "process invoice has no finishTime")
	}
}

func CreateInvoiceSimplifyForRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) *bean.Invoice {
	originalInvoice := query.GetInvoiceByInvoiceId(ctx, payment.InvoiceId)
	utility.Assert(originalInvoice != nil, "Payment Invoice Not found")
	var items []*bean.InvoiceItemSimplify
	err := utility.UnmarshalFromJsonString(originalInvoice.Lines, &items)
	if err != nil {
		return nil
	}
	var totalAmount int64 = 0
	for _, item := range items {
		totalAmount = totalAmount + item.Amount
	}
	var leftRefundAmount = refund.RefundAmount
	//proration to the items
	for _, item := range items {
		item.DiscountAmount = 0
		item.OriginAmount = item.Amount
		itemRefundAmount := utility.MinInt64(int64(float64(leftRefundAmount)*float64(item.Amount)/float64(totalAmount)), item.Amount)
		item.Amount = -itemRefundAmount
		item.Tax = int64(math.Floor(float64(item.Amount) * (1 - (1 / (1 + utility.ConvertTaxPercentageToInternalFloat(originalInvoice.TaxPercentage))))))
		item.AmountExcludingTax = item.Amount - item.Tax
		item.UnitAmountExcludingTax = int64(float64(item.AmountExcludingTax) / float64(item.Quantity))
		leftRefundAmount = leftRefundAmount - itemRefundAmount
	}
	//compensate to the items
	if leftRefundAmount > 0 {
		for _, item := range items {
			if leftRefundAmount > 0 {
				tempLeftDiscountAmount := utility.MinInt64(leftRefundAmount, item.OriginAmount-(-item.Amount))
				item.Amount = item.Amount - tempLeftDiscountAmount
				item.Tax = int64(math.Floor(float64(item.Amount) * (1 - (1 / (1 + utility.ConvertTaxPercentageToInternalFloat(originalInvoice.TaxPercentage))))))
				item.AmountExcludingTax = item.Amount - item.Tax
				item.UnitAmountExcludingTax = int64(float64(item.AmountExcludingTax) / float64(item.Quantity))
				leftRefundAmount = leftRefundAmount - tempLeftDiscountAmount
			} else {
				break
			}
		}
	}
	var refundType = "Partial Refund"
	if payment.TotalAmount == refund.RefundAmount {
		refundType = "Full Refund"
	}
	totalTax := -int64(math.Ceil(float64(refund.RefundAmount) * (1 - (1 / (1 + utility.ConvertTaxPercentageToInternalFloat(originalInvoice.TaxPercentage))))))
	return &bean.Invoice{
		InvoiceName:                    "Credit Note",
		ProductName:                    originalInvoice.ProductName,
		BizType:                        originalInvoice.BizType,
		Currency:                       originalInvoice.Currency,
		TaxAmount:                      totalTax,
		OriginAmount:                   -refund.RefundAmount,
		TotalAmount:                    -refund.RefundAmount,
		TotalAmountExcludingTax:        -refund.RefundAmount - totalTax,
		SubscriptionAmount:             -refund.RefundAmount,
		SubscriptionAmountExcludingTax: -refund.RefundAmount - totalTax,
		CountryCode:                    originalInvoice.CountryCode,
		VatNumber:                      originalInvoice.VatNumber,
		TaxPercentage:                  originalInvoice.TaxPercentage,
		DiscountAmount:                 0,
		SendStatus:                     consts.InvoiceSendStatusUnSend,
		DayUtilDue:                     consts.DEFAULT_DAY_UTIL_DUE,
		Lines:                          items,
		SendNote:                       fmt.Sprintf("%s (%s)", originalInvoice.InvoiceId, refundType),
		PaymentId:                      payment.PaymentId,
		RefundId:                       refund.RefundId,
	}
}

func ProrationDiscountToItem(totalDiscountAmount int64, totalTaxAmount int64, items []*bean.InvoiceItemSimplify) {
	if len(items) == 0 {
		fmt.Printf("ProrationDiscountToItem error: items is blank")
		return
	}
	if totalDiscountAmount <= 0 {
		return
	}
	var totalAmountExcludingTax int64 = 0
	for _, item := range items {
		totalAmountExcludingTax = totalAmountExcludingTax + item.AmountExcludingTax
	}
	if totalDiscountAmount > totalAmountExcludingTax {
		fmt.Printf("ProrationDiscountToItem error: totalDiscountAmount > totalAmountExcludingTax")
		return
	}
	var leftDiscountAmount = totalDiscountAmount
	var leftTotalTaxAmount = totalTaxAmount
	for _, item := range items {
		appendDiscountAmount := utility.MinInt64(int64(float64(totalDiscountAmount)*float64(item.AmountExcludingTax)/float64(totalAmountExcludingTax)), item.AmountExcludingTax)
		leftDiscountAmount = leftDiscountAmount - appendDiscountAmount
		item.DiscountAmount = item.DiscountAmount + appendDiscountAmount
		item.Tax = int64(float64(item.AmountExcludingTax-item.DiscountAmount) * utility.ConvertTaxPercentageToInternalFloat(item.TaxPercentage))
		item.Amount = item.AmountExcludingTax - item.DiscountAmount + item.Tax
		item.OriginAmount = item.Amount + item.DiscountAmount
	}
	//compensate to the first one
	if leftDiscountAmount > 0 {
		for _, item := range items {
			if leftDiscountAmount > 0 {
				appendDiscountAmount := utility.MinInt64(leftDiscountAmount, item.AmountExcludingTax)
				leftDiscountAmount = leftDiscountAmount - appendDiscountAmount
				item.DiscountAmount = item.DiscountAmount + appendDiscountAmount
				item.Tax = int64(float64(item.AmountExcludingTax-item.DiscountAmount) * utility.ConvertTaxPercentageToInternalFloat(item.TaxPercentage))
				item.Amount = item.AmountExcludingTax - item.DiscountAmount + item.Tax
				item.OriginAmount = item.Amount + item.DiscountAmount
			} else {
				break
			}
		}
	}
	for _, item := range items {
		leftTotalTaxAmount = leftTotalTaxAmount - item.Tax
	}
	//compensate to the first one
	if leftTotalTaxAmount != 0 {
		for _, item := range items {
			if leftTotalTaxAmount != 0 {
				item.Tax = item.Tax + leftTotalTaxAmount
				item.Amount = item.AmountExcludingTax + item.Tax
				item.OriginAmount = item.Amount + item.DiscountAmount
				leftTotalTaxAmount = 0
			} else {
				break
			}
		}
	}
}

func ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx context.Context, req *CalculateInvoiceReq) *bean.Invoice {
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
	var planAmountExcludingTax = req.Quantity * plan.Amount
	var planTaxAmount = int64(float64(planAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
	var name = plan.PlanName
	var description = fmt.Sprintf("%d * %s %s", req.Quantity, plan.PlanName, period)
	if req.ProductData != nil && len(req.ProductData.Name) > 0 {
		name = req.ProductData.Name
		if len(req.ProductData.Description) > 0 {
			description = req.ProductData.Description
		}
	}
	invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
		Currency:               req.Currency,
		OriginAmount:           planAmountExcludingTax + planTaxAmount,
		Amount:                 planAmountExcludingTax + planTaxAmount,
		Tax:                    planTaxAmount,
		TaxPercentage:          req.TaxPercentage,
		AmountExcludingTax:     planAmountExcludingTax,
		UnitAmountExcludingTax: plan.Amount,
		Quantity:               req.Quantity,
		Name:                   name,
		Description:            description,
		PdfDescription:         fmt.Sprintf("%d * %s %s", req.Quantity, plan.PlanName, period),
		Plan:                   bean.SimplifyPlan(plan),
	})
	for _, addon := range addons {
		var addonAmountExcludingTax = addon.Quantity * addon.AddonPlan.Amount
		var addonTaxAmount = int64(float64(addonAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			OriginAmount:           addonAmountExcludingTax + addonTaxAmount,
			Amount:                 addonAmountExcludingTax + addonTaxAmount,
			Tax:                    addonTaxAmount,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     addonAmountExcludingTax,
			UnitAmountExcludingTax: addon.AddonPlan.Amount,
			Quantity:               addon.Quantity,
			Name:                   addon.AddonPlan.PlanName,
			Description:            fmt.Sprintf("%d * %s %s", addon.Quantity, addon.AddonPlan.PlanName, period),
			Plan:                   addon.AddonPlan,
		})
	}

	discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, plan.MerchantId, totalAmountExcludingTax, req.Currency, req.DiscountCode, req.TimeNow), totalAmountExcludingTax)
	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
	ProrationDiscountToItem(discountAmount, taxAmount, invoiceItems)

	return &bean.Invoice{
		BizType:                        consts.BizTypeSubscription,
		InvoiceName:                    req.InvoiceName,
		ProductName:                    plan.PlanName,
		OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		DiscountAmount:                 discountAmount,
		DiscountCode:                   req.DiscountCode,
		TaxAmount:                      taxAmount,
		Currency:                       req.Currency,
		CountryCode:                    req.CountryCode,
		VatNumber:                      req.VatNumber,
		TaxPercentage:                  req.TaxPercentage,
		SubscriptionAmount:             totalAmountExcludingTax + discountAmount + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount,
		Lines:                          invoiceItems,
		PeriodStart:                    req.PeriodStart,
		PeriodEnd:                      req.PeriodEnd,
		FinishTime:                     req.FinishTime,
		SendStatus:                     consts.InvoiceSendStatusUnSend,
		DayUtilDue:                     3,
		BillingCycleAnchor:             req.BillingCycleAnchor,
		Metadata:                       req.Metadata,
		CreateFrom:                     req.CreateFrom,
	}
}

type ProrationPlanParam struct {
	PlanId   uint64
	Quantity int64
}

type CalculateProrationInvoiceReq struct {
	Currency           string                 `json:"currency"`
	DiscountCode       string                 `json:"discountCode"`
	TimeNow            int64                  `json:"TimeNow"`
	CountryCode        string                 `json:"countryCode"`
	VatNumber          string                 `json:"vatNumber"`
	TaxPercentage      int64                  `json:"taxPercentage"`
	ProrationDate      int64                  `json:"prorationStart"`
	PeriodStart        int64                  `json:"periodStart"`
	PeriodEnd          int64                  `json:"periodEnd"`
	FinishTime         int64                  `json:"finishTime"`
	OldProrationPlans  []*ProrationPlanParam  `json:"oldPlans"`
	NewProrationPlans  []*ProrationPlanParam  `json:"newPlans"`
	InvoiceName        string                 `json:"invoiceName"`
	ProductName        string                 `json:"productName"`
	BillingCycleAnchor int64                  `json:"billingCycleAnchor"             description:"billing_cycle_anchor"` // billing_cycle_anchor
	Metadata           map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

func ComputeSubscriptionProrationToFixedEndInvoiceDetailSimplify(ctx context.Context, req *CalculateProrationInvoiceReq) *bean.Invoice {
	if req.OldProrationPlans == nil {
		req.OldProrationPlans = make([]*ProrationPlanParam, 0)
	}
	if req.NewProrationPlans == nil {
		req.NewProrationPlans = make([]*ProrationPlanParam, 0)
	}
	newMap := make(map[uint64]*ProrationPlanParam)
	for _, planSub := range req.NewProrationPlans {
		newMap[planSub.PlanId] = planSub
	}

	utility.Assert(req.ProrationDate > 0, "Invalid ProrationDate")
	utility.Assert(req.PeriodStart <= req.ProrationDate && req.ProrationDate <= req.PeriodEnd, "System Error, Subscription Need Update")

	timeScale := int64((float64(req.PeriodEnd-req.ProrationDate) / float64(req.PeriodEnd-req.PeriodStart)) * 10000)
	var invoiceItems []*bean.InvoiceItemSimplify
	var totalAmountExcludingTax int64
	var merchantId uint64
	for _, oldPlanSub := range req.OldProrationPlans {
		plan := query.GetPlanById(ctx, oldPlanSub.PlanId)
		merchantId = plan.MerchantId
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(oldPlanSub.PlanId, 10))
		unitAmountExcludingTax := int64(float64(plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(timeScale))
		if newPlanSub, ok := newMap[oldPlanSub.PlanId]; ok {
			//new plan contain old
			quantityDiff := newPlanSub.Quantity - oldPlanSub.Quantity
			if quantityDiff > 0 {
				// quantity increase
				var amountExcludingTax = quantityDiff * unitAmountExcludingTax
				var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
				invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
					Currency:               req.Currency,
					OriginAmount:           amountExcludingTax + taxAmount,
					Amount:                 amountExcludingTax + taxAmount,
					Tax:                    taxAmount,
					TaxPercentage:          req.TaxPercentage,
					AmountExcludingTax:     amountExcludingTax,
					UnitAmountExcludingTax: unitAmountExcludingTax,
					Quantity:               quantityDiff,
					Name:                   plan.PlanName,
					Description:            fmt.Sprintf("Remaining Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.ProrationDate).Layout("2006-01-02")),
					Plan:                   bean.SimplifyPlan(plan),
				})
				totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
			} else if quantityDiff < 0 {
				// quantity decrease
				quantityDiff = -quantityDiff
				unitAmountExcludingTax = -unitAmountExcludingTax
				var amountExcludingTax = quantityDiff * unitAmountExcludingTax
				var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
				invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
					Currency:               req.Currency,
					OriginAmount:           amountExcludingTax + taxAmount,
					Amount:                 amountExcludingTax + taxAmount,
					Tax:                    taxAmount,
					TaxPercentage:          req.TaxPercentage,
					AmountExcludingTax:     amountExcludingTax,
					UnitAmountExcludingTax: unitAmountExcludingTax,
					Quantity:               quantityDiff,
					Name:                   plan.PlanName,
					Description:            fmt.Sprintf("Unused Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02")),
					Plan:                   bean.SimplifyPlan(plan),
				})
				totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
			}
			delete(newMap, newPlanSub.PlanId)
		} else {
			//old removed
			quantityDiff := oldPlanSub.Quantity
			unitAmountExcludingTax = -unitAmountExcludingTax
			var amountExcludingTax = quantityDiff * unitAmountExcludingTax
			var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
			invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
				Currency:               req.Currency,
				OriginAmount:           amountExcludingTax + taxAmount,
				Amount:                 amountExcludingTax + taxAmount,
				Tax:                    taxAmount,
				TaxPercentage:          req.TaxPercentage,
				AmountExcludingTax:     amountExcludingTax,
				UnitAmountExcludingTax: unitAmountExcludingTax,
				Quantity:               quantityDiff,
				Name:                   plan.PlanName,
				Description:            fmt.Sprintf("Unused Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02")),
				Plan:                   bean.SimplifyPlan(plan),
			})
			totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
		}
	}
	for _, newPlanSub := range newMap {
		plan := query.GetPlanById(ctx, newPlanSub.PlanId)
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(newPlanSub.PlanId, 10))
		unitAmountExcludingTax := int64(float64(plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(timeScale))
		quantityDiff := newPlanSub.Quantity
		var amountExcludingTax = quantityDiff * unitAmountExcludingTax
		var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			OriginAmount:           amountExcludingTax + taxAmount,
			Amount:                 amountExcludingTax + taxAmount,
			Tax:                    taxAmount,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: unitAmountExcludingTax,
			Quantity:               quantityDiff,
			Name:                   plan.PlanName,
			Description:            fmt.Sprintf("Remaining Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.ProrationDate).Layout("2006-01-02")),
			Plan:                   bean.SimplifyPlan(plan),
		})
		totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
	}

	discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, merchantId, totalAmountExcludingTax, req.Currency, req.DiscountCode, req.TimeNow), totalAmountExcludingTax)
	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
	ProrationDiscountToItem(discountAmount, taxAmount, invoiceItems)
	return &bean.Invoice{
		BizType:                        consts.BizTypeSubscription,
		InvoiceName:                    req.InvoiceName,
		ProductName:                    req.ProductName,
		OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		DiscountAmount:                 discountAmount,
		DiscountCode:                   req.DiscountCode,
		TaxAmount:                      taxAmount,
		Currency:                       req.Currency,
		CountryCode:                    req.CountryCode,
		VatNumber:                      req.VatNumber,
		TaxPercentage:                  req.TaxPercentage,
		SubscriptionAmount:             totalAmountExcludingTax + discountAmount + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount,
		Lines:                          invoiceItems,
		ProrationDate:                  req.ProrationDate,
		ProrationScale:                 timeScale,
		PeriodStart:                    req.ProrationDate,
		PeriodEnd:                      req.PeriodEnd,
		FinishTime:                     req.FinishTime,
		SendStatus:                     consts.InvoiceSendStatusUnSend,
		DayUtilDue:                     3,
		BillingCycleAnchor:             req.BillingCycleAnchor,
		Metadata:                       req.Metadata,
	}
}

func ComputeSubscriptionProrationToDifferentIntervalInvoiceDetailSimplify(ctx context.Context, req *CalculateProrationInvoiceReq) *bean.Invoice {
	if req.OldProrationPlans == nil {
		req.OldProrationPlans = make([]*ProrationPlanParam, 0)
	}
	if req.NewProrationPlans == nil {
		req.NewProrationPlans = make([]*ProrationPlanParam, 0)
	}
	newMap := make(map[uint64]*ProrationPlanParam)
	for _, planSub := range req.NewProrationPlans {
		newMap[planSub.PlanId] = planSub
	}

	utility.Assert(req.ProrationDate > 0, "Invalid ProrationDate")
	utility.Assert(req.PeriodStart <= req.ProrationDate && req.ProrationDate <= req.PeriodEnd, "System Error, Subscription Need Update")

	timeScale := int64((float64(req.PeriodEnd-req.ProrationDate) / float64(req.PeriodEnd-req.PeriodStart)) * 10000)
	var invoiceItems []*bean.InvoiceItemSimplify
	var totalAmountExcludingTax int64
	var merchantId uint64
	for _, oldPlanSub := range req.OldProrationPlans {
		plan := query.GetPlanById(ctx, oldPlanSub.PlanId)
		merchantId = plan.MerchantId
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(oldPlanSub.PlanId, 10))
		unitAmountExcludingTax := int64(float64(plan.Amount) * utility.ConvertTaxPercentageToInternalFloat(timeScale))
		//old removed
		quantityDiff := oldPlanSub.Quantity
		unitAmountExcludingTax = -unitAmountExcludingTax
		var amountExcludingTax = quantityDiff * unitAmountExcludingTax
		var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			OriginAmount:           amountExcludingTax + taxAmount,
			Amount:                 amountExcludingTax + taxAmount,
			Tax:                    taxAmount,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: unitAmountExcludingTax,
			Quantity:               quantityDiff,
			Name:                   plan.PlanName,
			Description:            fmt.Sprintf("Unused Time On %d * %s After %s", quantityDiff, plan.PlanName, gtime.NewFromTimeStamp(req.PeriodEnd).Layout("2006-01-02")),
			Plan:                   bean.SimplifyPlan(plan),
		})
		totalAmountExcludingTax = totalAmountExcludingTax + (quantityDiff * unitAmountExcludingTax)
	}
	var newPeriodEnd int64 = 0
	for _, newPlanSub := range newMap {
		plan := query.GetPlanById(ctx, newPlanSub.PlanId)
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(newPlanSub.PlanId, 10))
		if plan.Type == consts.PlanTypeMain {
			newPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, req.ProrationDate, req.ProrationDate, plan.Id)
		}
	}
	utility.Assert(newPeriodEnd > 0, "no main plan for upgrade")
	//change periodEnd
	req.PeriodEnd = newPeriodEnd

	for _, newPlanSub := range newMap {
		plan := query.GetPlanById(ctx, newPlanSub.PlanId)
		utility.Assert(plan != nil, "plan not found:"+strconv.FormatUint(newPlanSub.PlanId, 10))
		unitAmountExcludingTax := plan.Amount
		var amountExcludingTax = newPlanSub.Quantity * unitAmountExcludingTax
		var taxAmount = int64(float64(amountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
		invoiceItems = append(invoiceItems, &bean.InvoiceItemSimplify{
			Currency:               req.Currency,
			OriginAmount:           amountExcludingTax + taxAmount,
			Amount:                 amountExcludingTax + taxAmount,
			Tax:                    taxAmount,
			TaxPercentage:          req.TaxPercentage,
			AmountExcludingTax:     amountExcludingTax,
			UnitAmountExcludingTax: unitAmountExcludingTax,
			Quantity:               newPlanSub.Quantity,
			Name:                   plan.PlanName,
			Description:            fmt.Sprintf("%d * %s %s", newPlanSub.Quantity, plan.PlanName, gtime.NewFromTimeStamp(newPeriodEnd).Layout("2006-01-02")),
			Plan:                   bean.SimplifyPlan(plan),
		})
		totalAmountExcludingTax = totalAmountExcludingTax + (newPlanSub.Quantity * unitAmountExcludingTax)
	}

	utility.Assert(totalAmountExcludingTax >= 0, "not available for downgrade plan with different interval")
	discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, merchantId, totalAmountExcludingTax, req.Currency, req.DiscountCode, req.TimeNow), totalAmountExcludingTax)
	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(req.TaxPercentage))
	ProrationDiscountToItem(discountAmount, taxAmount, invoiceItems)
	return &bean.Invoice{
		BizType:                        consts.BizTypeSubscription,
		InvoiceName:                    req.InvoiceName,
		ProductName:                    req.ProductName,
		OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		DiscountAmount:                 discountAmount,
		DiscountCode:                   req.DiscountCode,
		TaxAmount:                      taxAmount,
		Currency:                       req.Currency,
		CountryCode:                    req.CountryCode,
		VatNumber:                      req.VatNumber,
		TaxPercentage:                  req.TaxPercentage,
		SubscriptionAmount:             totalAmountExcludingTax + discountAmount + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount,
		Lines:                          invoiceItems,
		ProrationDate:                  req.ProrationDate,
		ProrationScale:                 timeScale,
		PeriodStart:                    req.ProrationDate,
		PeriodEnd:                      req.PeriodEnd,
		FinishTime:                     req.FinishTime,
		SendStatus:                     consts.InvoiceSendStatusUnSend,
		DayUtilDue:                     3,
		BillingCycleAnchor:             req.BillingCycleAnchor,
		Metadata:                       req.Metadata,
	}
}
