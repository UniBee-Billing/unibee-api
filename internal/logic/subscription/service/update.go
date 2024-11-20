package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strconv"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/cmd/i18n"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/discount"
	service2 "unibee/internal/logic/gateway/service"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/service"
	plan2 "unibee/internal/logic/plan"
	subscription2 "unibee/internal/logic/subscription"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/subscription/pending_update_cancel"
	"unibee/internal/logic/user/sub_update"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func GetPlanIntervalLength(plan *entity.Plan) int {
	return plan2.PlanIntervalLength[plan.IntervalUnit] * plan.IntervalCount
}

func isUpgradeForSubscription(ctx context.Context, sub *entity.Subscription, plan *entity.Plan, quantity int64, addonParams []*bean.PlanAddonParam) (isUpgrade bool, changed bool) {
	//default logical，Effect Immediately for upgrade, effect at period end for downgrade
	//situation 1，NewPlan IntervalLength >  OldPlan IntervalLength，is upgrade，ignore Amount, Quantity and addon change
	//situation 2，NewPlan Unit Amount >  OldPlan Unit Amount，is upgrade，ignore Quantity and addon change
	//situation 3，NewPlan Unit Amount <  OldPlan Unit Amount，is downgrade，ignore Quantity and addon change
	//situation 4，NewPlan Total Amount >  OldPlan Total Amount，is upgrade
	//situation 5，NewPlan Total Amount <  OldPlan Total Amount，is downgrade
	//situation 6，NewPlan Total Amount =  OldPlan Total Amount，see Addon changes，if new addon appended or addon quantity changed, is upgrade，otherwise downgrade
	oldPlan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(oldPlan != nil, "old plan not found")
	utility.Assert(oldPlan.ProductId == plan.ProductId, "new plan's product must same with subscription")
	//if plan.IntervalUnit != oldPlan.IntervalUnit || plan.IntervalCount != oldPlan.IntervalCount {
	if GetPlanIntervalLength(plan) > GetPlanIntervalLength(oldPlan) {
		isUpgrade = true
		changed = true
	} else if plan.Amount > oldPlan.Amount || plan.Amount*quantity > oldPlan.Amount*sub.Quantity {
		isUpgrade = true
		changed = true
	} else if plan.Amount < oldPlan.Amount || plan.Amount*quantity < oldPlan.Amount*sub.Quantity {
		isUpgrade = false
		changed = true
	} else {
		var oldAddonParams []*bean.PlanAddonParam
		err := utility.UnmarshalFromJsonString(sub.AddonData, &oldAddonParams)
		utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString err:%v", err))
		var oldAddonMap = make(map[uint64]int64)
		for _, oldAddon := range oldAddonParams {
			if _, ok := oldAddonMap[oldAddon.AddonPlanId]; ok {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddonMap[oldAddon.AddonPlanId] + oldAddon.Quantity
			} else {
				oldAddonMap[oldAddon.AddonPlanId] = oldAddon.Quantity
			}
		}
		var newAddonMap = make(map[uint64]int64)
		for _, newAddon := range addonParams {
			if _, ok := newAddonMap[newAddon.AddonPlanId]; ok {
				newAddonMap[newAddon.AddonPlanId] = newAddonMap[newAddon.AddonPlanId] + newAddon.Quantity
			} else {
				newAddonMap[newAddon.AddonPlanId] = newAddon.Quantity
			}
		}
		for newAddonPlanId, newAddonQuantity := range newAddonMap {
			if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
				if oldAddonQuantity < newAddonQuantity {
					isUpgrade = true
					changed = true
					break
				}
			} else {
				isUpgrade = true
				changed = true
				break
			}
		}
		if len(oldAddonMap) != len(newAddonMap) {
			changed = true
		} else {
			for newAddonPlanId, newAddonQuantity := range newAddonMap {
				if oldAddonQuantity, ok := oldAddonMap[newAddonPlanId]; ok {
					if oldAddonQuantity != newAddonQuantity {
						changed = true
						break
					}
				} else {
					changed = true
					break
				}
			}
		}
	}
	return
}

type UpdatePreviewInternalReq struct {
	SubscriptionId  string                 `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId       uint64                 `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity        int64                  `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId       *uint64                `json:"gatewayId" dc:"Id" `
	EffectImmediate int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
	TaxPercentage   *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	ProductData     *bean.PlanProductParam `json:"productData"  dc:"ProductData"  `
	PaymentMethodId string
	Metadata        map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type UpdatePreviewInternalRes struct {
	Subscription          *entity.Subscription       `json:"subscription"`
	Plan                  *entity.Plan               `json:"plan"`
	Quantity              int64                      `json:"quantity"`
	Gateway               *entity.MerchantGateway    `json:"gateway"`
	MerchantInfo          *entity.Merchant           `json:"merchantInfo"`
	AddonParams           []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                []*bean.PlanAddonDetail    `json:"addons"`
	OriginAmount          int64                      `json:"originAmount"                `
	TotalAmount           int64                      `json:"totalAmount"`
	DiscountAmount        int64                      `json:"discountAmount"`
	Currency              string                     `json:"currency"`
	UserId                uint64                     `json:"userId"`
	OldPlan               *entity.Plan               `json:"oldPlan"`
	Invoice               *bean.Invoice              `json:"invoice"`
	NextPeriodInvoice     *bean.Invoice              `json:"nextPeriodInvoice"`
	ProrationDate         int64                      `json:"prorationDate"`
	EffectImmediate       bool                       `json:"EffectImmediate"`
	Gateways              []*bean.Gateway            `json:"gateways"`
	Changed               bool                       `json:"changed"`
	IsUpgrade             bool                       `json:"isUpgrade"`
	TaxPercentage         int64                      `json:"taxPercentage" `
	RecurringDiscountCode string                     `json:"recurringDiscountCode" `
	Discount              *bean.MerchantDiscountCode `json:"discount" `
	PaymentMethodId       string
}

func SubscriptionUpdatePreview(ctx context.Context, req *UpdatePreviewInternalReq, prorationDate int64, merchantMemberId int64) (res *UpdatePreviewInternalRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	// todo mark addon binding check

	plan := query.GetPlanById(ctx, req.NewPlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	utility.Assert(plan.Type == consts.PlanTypeMain, fmt.Sprintf("Plan Id:%v Not Main Type", plan.Id))
	gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, req.GatewayId, req.PaymentMethodId, sub.SubscriptionId)
	utility.Assert(gatewayId > 0, "gateway need specified")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(service2.IsGatewaySupportCountryCode(ctx, gateway, sub.CountryCode), "gateway not support")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	user := query.GetUserAccountById(ctx, sub.UserId)
	utility.Assert(user != nil, "user not found")
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)
	var subscriptionTaxPercentage = sub.TaxPercentage
	percentage, countryCode, vatNumber, err := sub_update.GetUserTaxPercentage(ctx, sub.UserId)
	if err == nil {
		subscriptionTaxPercentage = percentage
	}
	if req.TaxPercentage != nil {
		subscriptionTaxPercentage = *req.TaxPercentage
	}

	var currency = sub.Currency
	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Active status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanTypeRecurringAddon, fmt.Sprintf("Addon Id:%v Not Recurring Type", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
	}
	oldPlan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")

	var hasIntervalChange = false
	if req.NewPlanId != sub.PlanId {
		//utility.Assert(oldPlan.IntervalUnit == plan.IntervalUnit, "newPlan must have same recurring interval to old")
		//utility.Assert(oldPlan.IntervalCount == plan.IntervalCount, "newPlan must have same recurring interval to old")
		if oldPlan.IntervalCount != plan.IntervalCount || oldPlan.IntervalUnit != plan.IntervalUnit {
			hasIntervalChange = true
		}
	}

	var effectImmediate = false

	isUpgrade, changed := isUpgradeForSubscription(ctx, sub, plan, req.Quantity, req.AddonParams)
	utility.Assert(changed, "subscription update should have plan or addons changed")
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["SubscriptionUpdate"] = true
	if isUpgrade {
		req.Metadata["IsUpgrade"] = true
		effectImmediate = true
	} else {
		req.Metadata["IsUpgrade"] = false
		effectImmediate = config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).DowngradeEffectImmediately
	}

	if req.EffectImmediate > 0 {
		utility.Assert(req.EffectImmediate == 1 || req.EffectImmediate == 2, "EffectImmediate should be 1 or 2")
		if req.EffectImmediate == 1 {
			effectImmediate = true
		} else {
			effectImmediate = false
		}
	}

	//if !effectImmediate {
	//	utility.Assert(sub.CancelAtPeriodEnd == 0, "subscription will cancel at period end, should resume subscription first")
	//}

	var currentInvoice *bean.Invoice
	var nextPeriodInvoice *bean.Invoice
	var recurringDiscountCode string
	if prorationDate == 0 {
		prorationDate = time.Now().Unix()
		if sub.TestClock > prorationDate && !config2.GetConfigInstance().IsProd() {
			prorationDate = sub.TestClock
		}
	}

	if len(req.DiscountCode) > 0 {
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:     plan.MerchantId,
			UserId:         sub.UserId,
			DiscountCode:   req.DiscountCode,
			Currency:       sub.Currency,
			SubscriptionId: sub.SubscriptionId,
			PLanId:         req.NewPlanId,
			TimeNow:        utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
		})
		utility.Assert(canApply, message)
		if isRecurring {
			recurringDiscountCode = req.DiscountCode
		}
		//} else if len(sub.DiscountCode) > 0 {
		//	canApply, isRecurring, _ := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
		//		MerchantId:     sub.MerchantId,
		//		UserId:         sub.UserId,
		//		DiscountCode:   sub.DiscountCode,
		//		Currency:       sub.Currency,
		//		SubscriptionId: sub.SubscriptionId,
		//		PLanId:         req.NewPlanId,
		//		TimeNow:        utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
		//	})
		//	if canApply && isRecurring {
		//		req.DiscountCode = sub.DiscountCode
		//		recurringDiscountCode = sub.DiscountCode
		//	}
	}

	if effectImmediate {
		if !config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).UpgradeProration {
			// without proration, just generate next cycle
			currentInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:        "SubscriptionUpdate",
				Currency:           sub.Currency,
				DiscountCode:       req.DiscountCode,
				TimeNow:            prorationDate,
				PlanId:             req.NewPlanId,
				Quantity:           req.Quantity,
				AddonJsonData:      utility.MarshalToJsonString(req.AddonParams),
				CountryCode:        countryCode,
				VatNumber:          vatNumber,
				TaxPercentage:      subscriptionTaxPercentage,
				PeriodStart:        prorationDate,
				PeriodEnd:          subscription2.GetPeriodEndFromStart(ctx, prorationDate, prorationDate, req.NewPlanId),
				FinishTime:         prorationDate,
				ProductData:        req.ProductData,
				BillingCycleAnchor: prorationDate,
				Metadata:           req.Metadata,
			})
		} else if prorationDate < sub.CurrentPeriodStart {
			// after period end before trial end, also or sub data not sync or use testClock in stage env
			currentInvoice = &bean.Invoice{
				InvoiceName:                    "SubscriptionUpdate",
				ProductName:                    plan.PlanName,
				OriginAmount:                   0,
				TotalAmount:                    0,
				TotalAmountExcludingTax:        0,
				DiscountCode:                   req.DiscountCode,
				DiscountAmount:                 0,
				Currency:                       sub.Currency,
				TaxAmount:                      0,
				SubscriptionAmount:             0,
				SubscriptionAmountExcludingTax: 0,
				Lines:                          make([]*bean.InvoiceItemSimplify, 0),
				ProrationDate:                  prorationDate,
				PeriodStart:                    sub.CurrentPeriodStart,
				PeriodEnd:                      sub.CurrentPeriodEnd,
				Metadata:                       req.Metadata,
				CountryCode:                    countryCode,
				VatNumber:                      vatNumber,
				TaxPercentage:                  subscriptionTaxPercentage,
			}
		} else if prorationDate > sub.CurrentPeriodEnd {
			// after periodEnd, is not a currentInvoice, just use it
			currentInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:        "SubscriptionUpdate",
				Currency:           sub.Currency,
				DiscountCode:       req.DiscountCode,
				TimeNow:            prorationDate,
				PlanId:             req.NewPlanId,
				Quantity:           req.Quantity,
				AddonJsonData:      utility.MarshalToJsonString(req.AddonParams),
				CountryCode:        countryCode,
				VatNumber:          vatNumber,
				TaxPercentage:      subscriptionTaxPercentage,
				PeriodStart:        prorationDate,
				PeriodEnd:          subscription2.GetPeriodEndFromStart(ctx, prorationDate, prorationDate, req.NewPlanId),
				FinishTime:         prorationDate,
				ProductData:        req.ProductData,
				BillingCycleAnchor: prorationDate,
				Metadata:           req.Metadata,
			})
		} else {
			// currentInvoice
			var oldAddonParams []*bean.PlanAddonParam
			err = utility.UnmarshalFromJsonString(sub.AddonData, &oldAddonParams)
			utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
			var oldProrationPlanParams []*invoice_compute.ProrationPlanParam
			oldProrationPlanParams = append(oldProrationPlanParams, &invoice_compute.ProrationPlanParam{
				PlanId:   sub.PlanId,
				Quantity: sub.Quantity,
			})
			for _, addonParam := range oldAddonParams {
				oldProrationPlanParams = append(oldProrationPlanParams, &invoice_compute.ProrationPlanParam{
					PlanId:   addonParam.AddonPlanId,
					Quantity: addonParam.Quantity,
				})
			}
			var newProrationPlanParams []*invoice_compute.ProrationPlanParam
			newProrationPlanParams = append(newProrationPlanParams, &invoice_compute.ProrationPlanParam{
				PlanId:   req.NewPlanId,
				Quantity: req.Quantity,
			})
			for _, addonParam := range req.AddonParams {
				newProrationPlanParams = append(newProrationPlanParams, &invoice_compute.ProrationPlanParam{
					PlanId:   addonParam.AddonPlanId,
					Quantity: addonParam.Quantity,
				})
			}
			oldCode := ""
			latestPaidInvoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
			if latestPaidInvoice.Status == consts.InvoiceStatusPaid {
				oldCode = latestPaidInvoice.DiscountCode
			} else {
				latestPaidInvoice = query.GetSubLatestPaidInvoice(ctx, sub.SubscriptionId)
				if latestPaidInvoice != nil {
					oldCode = latestPaidInvoice.DiscountCode
				}
			}
			if !hasIntervalChange {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToFixedEndInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:       "SubscriptionUpdate",
					ProductName:       plan.PlanName,
					Currency:          sub.Currency,
					DiscountCode:      req.DiscountCode,
					TimeNow:           prorationDate,
					CountryCode:       countryCode,
					VatNumber:         vatNumber,
					TaxPercentage:     subscriptionTaxPercentage,
					ProrationDate:     prorationDate,
					OldProrationPlans: oldProrationPlanParams,
					NewProrationPlans: newProrationPlanParams,
					PeriodStart:       sub.CurrentPeriodStart,
					PeriodEnd:         sub.CurrentPeriodEnd,
					Metadata:          req.Metadata,
					OldDiscountCode:   oldCode,
					OldTaxPercentage:  sub.TaxPercentage,
				})
			} else {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToDifferentIntervalInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:        "SubscriptionUpdate",
					ProductName:        plan.PlanName,
					Currency:           sub.Currency,
					DiscountCode:       req.DiscountCode,
					TimeNow:            prorationDate,
					CountryCode:        countryCode,
					VatNumber:          vatNumber,
					TaxPercentage:      subscriptionTaxPercentage,
					ProrationDate:      prorationDate,
					OldProrationPlans:  oldProrationPlanParams,
					NewProrationPlans:  newProrationPlanParams,
					PeriodStart:        sub.CurrentPeriodStart,
					PeriodEnd:          sub.CurrentPeriodEnd,
					BillingCycleAnchor: prorationDate,
					Metadata:           req.Metadata,
					OldDiscountCode:    oldCode,
					OldTaxPercentage:   sub.TaxPercentage,
				})
			}
		}
		prorationDate = currentInvoice.ProrationDate
	} else {
		prorationDate = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
		currentInvoice = &bean.Invoice{
			InvoiceName:                    "SubscriptionUpdate",
			ProductName:                    plan.PlanName,
			OriginAmount:                   0,
			TotalAmount:                    0,
			TotalAmountExcludingTax:        0,
			DiscountCode:                   req.DiscountCode,
			DiscountAmount:                 0,
			Currency:                       currency,
			TaxAmount:                      0,
			SubscriptionAmount:             0,
			SubscriptionAmountExcludingTax: 0,
			Lines:                          make([]*bean.InvoiceItemSimplify, 0),
			ProrationDate:                  prorationDate,
			PeriodStart:                    sub.CurrentPeriodStart,
			PeriodEnd:                      sub.CurrentPeriodEnd,
			Metadata:                       req.Metadata,
			CountryCode:                    countryCode,
			VatNumber:                      vatNumber,
			TaxPercentage:                  subscriptionTaxPercentage,
		}
	}

	nextCode := ""
	if len(recurringDiscountCode) > 0 {
		code := query.GetDiscountByCode(ctx, sub.MerchantId, recurringDiscountCode)
		if code.CycleLimit > 1 || code.CycleLimit == 0 {
			nextCode = recurringDiscountCode
		}
	}
	nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:        "SubscriptionCycle",
		Currency:           sub.Currency,
		DiscountCode:       nextCode,
		TimeNow:            prorationDate,
		PlanId:             req.NewPlanId,
		Quantity:           req.Quantity,
		AddonJsonData:      utility.MarshalToJsonString(req.AddonParams),
		CountryCode:        countryCode,
		VatNumber:          vatNumber,
		TaxPercentage:      subscriptionTaxPercentage,
		PeriodStart:        utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd),
		PeriodEnd:          subscription2.GetPeriodEndFromStart(ctx, utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd), prorationDate, req.NewPlanId),
		FinishTime:         utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd),
		ProductData:        req.ProductData,
		BillingCycleAnchor: prorationDate,
		Metadata:           req.Metadata,
	})

	if currentInvoice.TotalAmount <= 0 {
		effectImmediate = config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).DowngradeEffectImmediately
	}

	return &UpdatePreviewInternalRes{
		Subscription:          sub,
		Plan:                  plan,
		Quantity:              req.Quantity,
		Gateway:               gateway,
		MerchantInfo:          merchantInfo,
		AddonParams:           req.AddonParams,
		Addons:                addons,
		Currency:              currency,
		UserId:                sub.UserId,
		OldPlan:               oldPlan,
		OriginAmount:          currentInvoice.OriginAmount,
		TotalAmount:           currentInvoice.TotalAmount,
		DiscountAmount:        currentInvoice.DiscountAmount,
		Invoice:               currentInvoice,
		NextPeriodInvoice:     nextPeriodInvoice,
		ProrationDate:         prorationDate,
		EffectImmediate:       effectImmediate,
		Gateways:              service2.GetMerchantAvailableGatewaysByCountryCode(ctx, sub.MerchantId, sub.CountryCode),
		Changed:               changed,
		IsUpgrade:             isUpgrade,
		TaxPercentage:         subscriptionTaxPercentage,
		RecurringDiscountCode: recurringDiscountCode,
		Discount:              bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, currentInvoice.DiscountCode)),
		PaymentMethodId:       paymentMethodId,
	}, nil
}

type UpdateInternalReq struct {
	SubscriptionId     string                      `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId          uint64                      `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          *uint64                     `json:"gatewayId" dc:"GatewayId" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                      `json:"confirmCurrency" dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ProrationDate      *int64                      `json:"prorationDate" dc:"The utc time to start Proration, default current time" `
	EffectImmediate    int                         `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	Metadata           map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	DiscountCode       string                      `json:"discountCode"        dc:"DiscountCode"`
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment      bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl          string                      `json:"returnUrl"  dc:"ReturnUrl"  `
	CancelUrl          string                      `json:"cancelUrl" dc:"CancelUrl"`
	ProductData        *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
}

type UpdateInternalRes struct {
	SubscriptionPendingUpdate *detail.SubscriptionPendingUpdateDetail `json:"subscriptionPendingUpdate" dc:"SubscriptionPendingUpdate"`
	Paid                      bool                                    `json:"paid" dc:"Paid，true|false"`
	Link                      string                                  `json:"link" dc:"Pay Link"`
	Note                      string                                  `json:"note" dc:"note"`
}

func SubscriptionUpdate(ctx context.Context, req *UpdateInternalReq, merchantMemberId int64) (*UpdateInternalRes, error) {
	var prorationDate int64 = 0
	if req.ProrationDate != nil {
		prorationDate = *req.ProrationDate
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	if req.Discount != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		// create external discount
		utility.Assert(req.NewPlanId > 0, "planId invalid")
		utility.Assert(sub.UserId > 0, "UserId invalid")
		plan := query.GetPlanById(ctx, req.NewPlanId)
		utility.Assert(plan.MerchantId == sub.MerchantId, "merchant not match")
		utility.Assert(plan != nil, "invalid planId")
		one := discount.CreateExternalDiscount(ctx, sub.MerchantId, sub.UserId, strconv.FormatUint(req.NewPlanId, 10), req.Discount, plan.Currency, utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock))
		req.DiscountCode = one.Code
	} else if len(req.DiscountCode) > 0 {
		one := query.GetDiscountByCode(ctx, sub.MerchantId, req.DiscountCode)
		utility.Assert(one.Type == 0, "invalid code, code is from external")
	}
	prepare, err := SubscriptionUpdatePreview(ctx, &UpdatePreviewInternalReq{
		SubscriptionId:  req.SubscriptionId,
		NewPlanId:       req.NewPlanId,
		Quantity:        req.Quantity,
		AddonParams:     req.AddonParams,
		GatewayId:       req.GatewayId,
		EffectImmediate: req.EffectImmediate,
		DiscountCode:    req.DiscountCode,
		TaxPercentage:   req.TaxPercentage,
		ProductData:     req.ProductData,
		Metadata:        req.Metadata,
	}, prorationDate, merchantMemberId)
	if err != nil {
		return nil, err
	}

	//subscription prepare
	if req.ConfirmTotalAmount > 0 {
		utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, i18n.LocalizationFormat(ctx, "{#AmountNotMatch}"))
	}
	if len(req.ConfirmCurrency) > 0 {
		utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	}
	if prepare.Invoice.TotalAmount <= 0 {
		utility.Assert(prepare.EffectImmediate == config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).DowngradeEffectImmediately, "System Error, Cannot Effect Immediate With Negative CaptureAmount")
	}

	var effectImmediate = 0
	var effectTime = prepare.Subscription.CurrentPeriodEnd
	if prepare.EffectImmediate {
		effectImmediate = 1
		effectTime = gtime.Now().Timestamp()
	}

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:       prepare.MerchantInfo.Id,
		GatewayId:        prepare.Gateway.Id,
		UserId:           prepare.Subscription.UserId,
		SubscriptionId:   prepare.Subscription.SubscriptionId,
		PendingUpdateId:  utility.CreatePendingUpdateId(),
		Amount:           prepare.Subscription.Amount,
		Currency:         prepare.Subscription.Currency,
		PlanId:           prepare.Subscription.PlanId,
		Quantity:         prepare.Subscription.Quantity,
		AddonData:        prepare.Subscription.AddonData,
		UpdateAmount:     prepare.NextPeriodInvoice.TotalAmount,
		ProrationAmount:  prepare.Invoice.TotalAmount,
		UpdateCurrency:   prepare.Currency,
		UpdatePlanId:     prepare.Plan.Id,
		UpdateQuantity:   prepare.Quantity,
		UpdateAddonData:  utility.MarshalToJsonString(prepare.AddonParams),
		Status:           consts.PendingSubStatusInit,
		Data:             "",
		MerchantMemberId: merchantMemberId,
		ProrationDate:    prorationDate,
		EffectImmediate:  effectImmediate,
		EffectTime:       effectTime,
		TaxPercentage:    prepare.TaxPercentage,
		DiscountCode:     prepare.RecurringDiscountCode,
		CreateTime:       gtime.Now().Timestamp(),
		MetaData:         utility.MarshalToJsonString(prepare.Invoice.Metadata),
	}

	result, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	if prepare.Invoice.Metadata == nil {
		prepare.Invoice.Metadata = make(map[string]interface{})
	}
	prepare.Invoice.Metadata["SubscriptionPendingUpdateId"] = one.PendingUpdateId
	var subUpdateRes *UpdateSubscriptionInternalResp
	if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
		// createAndPayNewProrationInvoice
		merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
		utility.Assert(merchantInfo != nil, "merchantInfo not found")
		// utility.Assert(user != nil, "user not found")
		invoice, err := service3.CreateProcessingInvoiceForSub(ctx, req.NewPlanId, prepare.Invoice, prepare.Subscription, prepare.Gateway.Id, prepare.PaymentMethodId, false, prepare.ProrationDate)
		utility.AssertError(err, "System Error")
		createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, invoice, req.ManualPayment, req.ReturnUrl, req.CancelUrl, "SubscriptionUpdate", 0)
		if err != nil {
			g.Log().Print(ctx, "SubscriptionUpdate CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return nil, err
		}
		// Upgrade
		subUpdateRes = &UpdateSubscriptionInternalResp{
			GatewayUpdateId: invoice.InvoiceId,
			Data:            utility.MarshalToJsonString(createRes),
			Link:            createRes.Link,
			Paid:            createRes.Status == consts.PaymentSuccess,
			Invoice:         createRes.Invoice,
		}
	} else if prepare.EffectImmediate && prepare.Invoice.TotalAmount == 0 {
		//totalAmount is 0, no payment need
		invoice, err := service3.CreateProcessingInvoiceForSub(ctx, req.NewPlanId, prepare.Invoice, prepare.Subscription, prepare.Gateway.Id, prepare.PaymentMethodId, false, prepare.ProrationDate)
		utility.AssertError(err, "System Error")
		invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
		utility.AssertError(err, "System Error")
		subUpdateRes = &UpdateSubscriptionInternalResp{
			GatewayUpdateId: invoice.InvoiceId,
			Paid:            true,
			Link:            GetSubscriptionZeroPaymentLink(req.ReturnUrl, sub.SubscriptionId),
			Invoice:         invoice,
		}
	} else {
		prepare.EffectImmediate = false
		effectImmediate = 0
		subUpdateRes = &UpdateSubscriptionInternalResp{
			Paid: false,
			Link: "",
		}
	}

	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
		Content:        "Update",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      one.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)

	if err != nil {
		return nil, err
	}
	// bing to subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().PendingUpdateId: one.PendingUpdateId,
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	// only one need, cancel others
	// need cancel payment、 invoice and send invoice email
	pending_update_cancel.CancelOtherUnfinishedPendingUpdatesBackground(prepare.Subscription.SubscriptionId, one.PendingUpdateId, "CancelByNewUpdate-"+one.PendingUpdateId)

	var PaidInt = 0
	if subUpdateRes.Paid {
		PaidInt = 1
	}
	var note = "Success"
	if effectImmediate == 1 && !subUpdateRes.Paid {
		note = "Payment Action Required"
	} else if effectImmediate == 0 {
		note = "Will Effect At Period End"
	}

	one.Link = subUpdateRes.Link
	one.Status = consts.PendingSubStatusCreate
	_, err = dao.SubscriptionPendingUpdate.Ctx(ctx).Data(g.Map{
		dao.SubscriptionPendingUpdate.Columns().Status:          consts.PendingSubStatusCreate,
		dao.SubscriptionPendingUpdate.Columns().ResponseData:    subUpdateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:       gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Paid:            PaidInt,
		dao.SubscriptionPendingUpdate.Columns().Link:            subUpdateRes.Link,
		dao.SubscriptionPendingUpdate.Columns().InvoiceId:       subUpdateRes.GatewayUpdateId,
		dao.SubscriptionPendingUpdate.Columns().Note:            note,
		dao.SubscriptionPendingUpdate.Columns().MetaData:        utility.MarshalToJsonString(prepare.Invoice.Metadata),
		dao.SubscriptionPendingUpdate.Columns().EffectImmediate: effectImmediate,
	}).Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, one.PendingUpdateId).OmitNil().Update()
	if err != nil {
		return nil, err
	} else {
		_, _ = redismq.Send(&redismq.Message{
			Topic:      redismq2.TopicSubscriptionPendingUpdateCreate.Topic,
			Tag:        redismq2.TopicSubscriptionPendingUpdateCreate.Tag,
			Body:       one.PendingUpdateId,
			CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
		})
	}

	if prepare.EffectImmediate && subUpdateRes.Paid {
		_, err = handler.HandlePendingUpdatePaymentSuccess(ctx, prepare.Subscription, one.PendingUpdateId, subUpdateRes.Invoice)
		if err != nil {
			return nil, err
		}
		one.Status = consts.PendingSubStatusFinished
	}

	return &UpdateInternalRes{
		SubscriptionPendingUpdate: &detail.SubscriptionPendingUpdateDetail{
			MerchantId:      one.MerchantId,
			SubscriptionId:  one.SubscriptionId,
			PendingUpdateId: one.PendingUpdateId,
			GmtCreate:       one.GmtCreate,
			Amount:          one.Amount,
			Status:          one.Status,
			UpdateAmount:    one.UpdateAmount,
			Currency:        one.Currency,
			UpdateCurrency:  one.UpdateCurrency,
			PlanId:          one.PlanId,
			UpdatePlanId:    one.UpdatePlanId,
			Quantity:        one.Quantity,
			UpdateQuantity:  one.UpdateQuantity,
			AddonData:       one.AddonData,
			UpdateAddonData: one.UpdateAddonData,
			ProrationAmount: one.ProrationAmount,
			GatewayId:       one.GatewayId,
			UserId:          one.UserId,
			InvoiceId:       one.InvoiceId,
			GmtModify:       one.GmtModify,
			Paid:            one.Paid,
			Link:            one.Link,
			MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate: effectImmediate,
			EffectTime:      one.EffectTime,
			Note:            one.Note,
			Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:        prepare.Invoice.Metadata,
		},
		Paid: len(subUpdateRes.Link) == 0 || subUpdateRes.Paid, // link is blank or paid is true, portal will not redirect
		Link: subUpdateRes.Link,
		Note: note,
	}, nil
}

type UpdateSubscriptionInternalResp struct {
	GatewayUpdateId string          `json:"gatewayUpdateId" description:""`
	Data            string          `json:"data"`
	Link            string          `json:"link" description:""`
	Paid            bool            `json:"paid" description:""`
	Invoice         *entity.Invoice `json:"invoice" description:""`
}
