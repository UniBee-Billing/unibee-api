package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"strings"
	"time"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/user/subscription"
	"unibee/api/user/vat"
	config2 "unibee/internal/cmd/config"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/gateway_bean"
	service2 "unibee/internal/logic/gateway/service"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/payment/method"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

func checkAndListAddonsFromParams(ctx context.Context, addonParams []*bean.PlanAddonParam) []*bean.PlanAddonDetail {
	var addons []*bean.PlanAddonDetail
	var totalAddonIds []uint64
	if len(addonParams) > 0 {
		for _, s := range addonParams {
			totalAddonIds = append(totalAddonIds, s.AddonPlanId) // 添加到整数列表中
		}
	}
	var allAddonList []*entity.Plan
	if len(totalAddonIds) > 0 {
		//query all plan
		err := dao.Plan.Ctx(ctx).WhereIn(dao.Plan.Columns().Id, totalAddonIds).OmitEmpty().Scan(&allAddonList)
		if err == nil {
			//add to list
			mapPlans := make(map[uint64]*entity.Plan)
			for _, pair := range allAddonList {
				key := pair.Id
				value := pair
				mapPlans[key] = value
			}
			for _, param := range addonParams {
				utility.Assert(mapPlans[param.AddonPlanId] != nil, fmt.Sprintf("AddonPlanId not found:%v", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("Id:%v not Addon Type", param.AddonPlanId))
				utility.Assert(mapPlans[param.AddonPlanId].IsDeleted == 0, fmt.Sprintf("Addon Id:%v is Deleted", param.AddonPlanId))
				utility.Assert(param.Quantity > 0, fmt.Sprintf("Id:%v quantity invalid", param.AddonPlanId))
				addons = append(addons, &bean.PlanAddonDetail{
					Quantity:  param.Quantity,
					AddonPlan: bean.SimplifyPlan(mapPlans[param.AddonPlanId]),
				})
			}
		}
	}
	return addons
}

func VatNumberValidate(ctx context.Context, req *vat.NumberValidateReq, userId uint64) (*vat.NumberValidateRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), userId, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	if vatNumberValidate.Valid {
		vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, _interface.GetMerchantId(ctx), vatNumberValidate.CountryCode)
		utility.Assert(err == nil, fmt.Sprintf("verify error:%s", err))
		utility.Assert(vatCountryRate != nil, fmt.Sprintf("vatNumber not found for countryCode:%v", vatNumberValidate.CountryCode))
	}
	return &vat.NumberValidateRes{VatNumberValidate: vatNumberValidate}, nil
}

func MerchantGatewayCheck(ctx context.Context, merchantId uint64, reqGatewayId uint64) *entity.MerchantGateway {
	if reqGatewayId > 0 {
		gateway := query.GetGatewayById(ctx, reqGatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		utility.Assert(gateway.MerchantId == merchantId, "gateway not match")
		return gateway
	} else {
		list := query.GetMerchantGatewayList(ctx, merchantId)
		utility.Assert(len(list) > 0, "merchant gateway need setup")
		utility.Assert(len(list) == 1, "gateway need specify")
		return list[0]
	}
}

type RenewInternalReq struct {
	MerchantId     uint64 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	SubscriptionId string `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	//UserId         uint64                      `json:"userId" dc:"UserId" v:"required"`
	GatewayId     *uint64                     `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	TaxPercentage *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	DiscountCode  string                      `json:"discountCode" dc:"DiscountCode, override subscription discount"`
	Discount      *bean.ExternalDiscountParam `json:"discount" dc:"Discount, override subscription discount"`
	ManualPayment bool                        `json:"manualPayment" dc:"ManualPayment"`
	ReturnUrl     string                      `json:"returnUrl"  dc:"ReturnUrl"  `
}

func SubscriptionRenew(ctx context.Context, req *RenewInternalReq) (*CreateInternalRes, error) {
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.MerchantId == req.MerchantId, "merchantId not match")
	// todo mark renew for all status
	//utility.Assert(sub.Status == consts.SubStatusExpired || sub.Status == consts.SubStatusCancelled, "subscription not cancel or expire status")
	var subscriptionTaxPercentage = sub.TaxPercentage
	if req.TaxPercentage != nil {
		subscriptionTaxPercentage = *req.TaxPercentage
	}
	var addonParams []*bean.PlanAddonParam
	if len(sub.AddonData) > 0 {
		err := utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
		if err != nil {
			g.Log().Errorf(ctx, "SubscriptionDetail Unmarshal addon param:%s", err.Error())
		}
	}
	var gatewayId = sub.GatewayId
	if req.GatewayId != nil {
		gatewayId = *req.GatewayId
	}

	var timeNow = gtime.Now().Timestamp()
	if sub.TestClock > sub.CurrentPeriodStart && !config2.GetConfigInstance().IsProd() {
		timeNow = sub.TestClock
	}

	if req.Discount != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		// create external discount
		utility.Assert(sub.PlanId > 0, "planId invalid")
		plan := query.GetPlanById(ctx, sub.PlanId)
		utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
		utility.Assert(plan != nil, "invalid planId")
		one := discount.CreateExternalDiscount(ctx, req.MerchantId, sub.UserId, strconv.FormatUint(sub.PlanId, 10), req.Discount, plan.Currency)
		req.DiscountCode = one.Code
	} else if len(req.DiscountCode) > 0 {
		one := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
		utility.Assert(one.Type == 0, "invalid code, code is from external")
	}

	if len(req.DiscountCode) > 0 {
		canApply, _, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:     sub.MerchantId,
			UserId:         sub.UserId,
			DiscountCode:   req.DiscountCode,
			Currency:       sub.Currency,
			SubscriptionId: sub.SubscriptionId,
			PLanId:         sub.PlanId,
		})
		utility.Assert(canApply, message)
	} else if len(req.DiscountCode) == 0 && len(sub.DiscountCode) > 0 {
		canApply, isRecurring, _ := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:     sub.MerchantId,
			UserId:         sub.UserId,
			DiscountCode:   sub.DiscountCode,
			Currency:       sub.Currency,
			SubscriptionId: sub.SubscriptionId,
			PLanId:         sub.PlanId,
		})
		if canApply && isRecurring {
			req.DiscountCode = sub.DiscountCode
		}
	}

	currentInvoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:   "SubscriptionCycle",
		Currency:      sub.Currency,
		DiscountCode:  req.DiscountCode,
		TimeNow:       timeNow,
		PlanId:        sub.PlanId,
		Quantity:      sub.Quantity,
		AddonJsonData: utility.MarshalToJsonString(addonParams),
		TaxPercentage: subscriptionTaxPercentage,
		PeriodStart:   timeNow,
		PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, timeNow, sub.PlanId),
		FinishTime:    timeNow,
	})

	// createAndPayNewProrationInvoice
	merchantInfo := query.GetMerchantById(ctx, sub.MerchantId)
	utility.Assert(merchantInfo != nil, "merchantInfo not found")
	// utility.Assert(user != nil, "user not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, currentInvoice, sub)
	utility.AssertError(err, "System Error")
	createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, sub.GatewayDefaultPaymentMethod, invoice, gateway.Id, req.ManualPayment, req.ReturnUrl, "SubscriptionRenew")
	if err != nil {
		g.Log().Print(ctx, "SubscriptionRenew CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
		return nil, err
	}

	// need cancel payment、 invoice and send invoice email
	CancelOtherUnfinishedPendingUpdatesBackground(sub.SubscriptionId, sub.SubscriptionId, "CancelByRenewSubscription-"+sub.SubscriptionId)

	if createRes.Status == consts.PaymentSuccess && createRes.Payment != nil {
		err = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, createRes.Payment)
		if err != nil {
			return nil, err
		}
	}

	sub = query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)

	_, _ = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().TrialEnd:          sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()

	return &CreateInternalRes{
		Subscription: bean.SimplifySubscription(sub),
		Paid:         createRes.Status == consts.PaymentSuccess && createRes.Payment != nil,
		Link:         createRes.Link,
	}, nil
}

type CreatePreviewInternalReq struct {
	MerchantId     uint64                 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId         uint64                 `json:"planId" dc:"PlanId" v:"required"`
	UserId         uint64                 `json:"userId" dc:"UserId" v:"required"`
	Quantity       int64                  `json:"quantity" dc:"Quantity" `
	DiscountCode   string                 `json:"discountCode"        dc:"DiscountCode"`
	GatewayId      uint64                 `json:"gatewayId" dc:"Id" v:"required" `
	AddonParams    []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber      string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage  *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	TrialEnd       int64                  `json:"trialEnd"  description:"trial_end, utc time"` // trial_end, utc time
	IsSubmit       bool
}

type CreatePreviewInternalRes struct {
	Plan                     *entity.Plan                       `json:"plan"`
	Quantity                 int64                              `json:"quantity"`
	Gateway                  *entity.MerchantGateway            `json:"gateway"`
	Merchant                 *entity.Merchant                   `json:"merchantInfo"`
	AddonParams              []*bean.PlanAddonParam             `json:"addonParams"`
	Addons                   []*bean.PlanAddonDetail            `json:"addons"`
	OriginAmount             int64                              `json:"originAmount"                `
	TotalAmount              int64                              `json:"totalAmount" `
	DiscountAmount           int64                              `json:"discountAmount"`
	Currency                 string                             `json:"currency" `
	VatCountryCode           string                             `json:"vatCountryCode" `
	VatCountryName           string                             `json:"vatCountryName" `
	VatNumber                string                             `json:"vatNumber" `
	VatNumberValidate        *bean.ValidResult                  `json:"vatNumberValidate" `
	TaxPercentage            int64                              `json:"taxPercentage" `
	TrialEnd                 int64                              `json:"trialEnd" `
	VatVerifyData            string                             `json:"vatVerifyData" `
	Invoice                  *bean.InvoiceSimplify              `json:"invoice"`
	UserId                   uint64                             `json:"userId" `
	Email                    string                             `json:"email" `
	VatCountryRate           *bean.VatCountryRate               `json:"vatCountryRate" `
	Gateways                 []*bean.GatewaySimplify            `json:"gateways" `
	RecurringDiscountCode    string                             `json:"recurringDiscountCode" `
	Discount                 *bean.MerchantDiscountCodeSimplify `json:"discount" `
	VatNumberValidateMessage string                             `json:"vatNumberValidateMessage" `
	DiscountMessage          string                             `json:"discountMessage" `
	CancelAtPeriodEnd        int                                `json:"cancelAtPeriodEnd"           description:"whether cancel at period end，0-false | 1-true"` // whether cancel at period end，0-false | 1-true
}

type CreateInternalReq struct {
	MerchantId         uint64                      `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId             uint64                      `json:"planId" dc:"PlanId" v:"required"`
	UserId             uint64                      `json:"userId" dc:"UserId" v:"required"`
	DiscountCode       string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          uint64                      `json:"gatewayId" dc:"Id"   v:"required" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                      `json:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl          string                      `json:"returnUrl"  dc:"RedirectUrl"  `
	VatCountryCode     string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	PaymentMethodId    string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata           map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	TrialEnd           int64                       `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	StartIncomplete    bool                        `json:"StartIncomplete"        dc:"StartIncomplete, use now pay later, subscription will generate invoice and start with incomplete status if set"`
}

type CreateInternalRes struct {
	Subscription *bean.SubscriptionSimplify `json:"subscription" dc:"Subscription"`
	Paid         bool                       `json:"paid"`
	Link         string                     `json:"link"`
}

func SubscriptionCreatePreview(ctx context.Context, req *CreatePreviewInternalReq) (*CreatePreviewInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.GatewayId > 0, "GatewayId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	utility.Assert(plan.Type == consts.PlanTypeMain, fmt.Sprintf("Plan Id:%v Not Main Type", plan.Id))
	gateway := MerchantGatewayCheck(ctx, plan.MerchantId, req.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(service2.IsGatewaySupportCountryCode(ctx, gateway, req.VatCountryCode), "gateway not support")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	req.Quantity = utility.MaxInt64(1, req.Quantity)

	var err error
	utility.Assert(query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, req.UserId, merchantInfo.Id) == nil, "Another pending or active subscription exist")

	//vat
	if len(req.VatCountryCode) == 0 && len(user.CountryCode) > 0 {
		req.VatCountryCode = user.CountryCode
	}
	var vatCountryCode = req.VatCountryCode
	var subscriptionTaxPercentage int64 = 0
	var vatCountryName = ""
	var vatCountryRate *bean.VatCountryRate
	var vatNumberValidate *bean.ValidResult
	var vatNumberValidateMessage string
	var discountMessage string

	if len(req.VatNumber) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil, "Vat VATGateway need setup")
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.UserId, req.VatNumber, "")
		if err != nil || !vatNumberValidate.Valid {
			if err != nil {
				g.Log().Error(ctx, "ValidateVatNumberByDefaultGateway error:%s", err.Error())
				vatNumberValidateMessage = "Server Error"
			} else {
				vatNumberValidateMessage = "Validate Failure"
			}
		} else {
			vatCountryCode = vatNumberValidate.CountryCode
		}
		//if err != nil {
		//	return nil, err
		//}
		if req.IsSubmit {
			utility.Assert(vatNumberValidate.Valid, fmt.Sprintf("VatNumber validate failure, number:"+req.VatNumber))
		}
		//vatCountryCode = vatNumberValidate.CountryCode
	}

	if req.TaxPercentage != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "External TaxPercentage only available for api call")
		utility.Assert(*req.TaxPercentage > 0 && *req.TaxPercentage < 10000, "invalid taxPercentage")
		subscriptionTaxPercentage = *req.TaxPercentage
	} else if len(vatCountryCode) > 0 {
		if vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil {
			vatCountryRate, err = vat_gateway.QueryVatCountryRateByMerchant(ctx, merchantInfo.Id, vatCountryCode)
			if err == nil && vatCountryRate != nil {
				vatCountryName = vatCountryRate.CountryName
				if vatNumberValidate != nil && !strings.Contains(config2.GetConfigInstance().VatConfig.NumberUnExemptionCountryCodes, vatCountryCode) {
					subscriptionTaxPercentage = 0
				} else if vatCountryRate.StandardTaxPercentage > 0 {
					subscriptionTaxPercentage = vatCountryRate.StandardTaxPercentage
				}
			}
		}
	}

	var currency = plan.Currency
	var TotalAmountExcludingTax = plan.Amount * req.Quantity

	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)

	for _, addon := range addons {
		utility.Assert(strings.Compare(addon.AddonPlan.Currency, currency) == 0, fmt.Sprintf("currency not match for planId:%v addonId:%v", plan.Id, addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.MerchantId == plan.MerchantId, fmt.Sprintf("Addon Id:%v Merchant not match", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Status == consts.PlanStatusActive, fmt.Sprintf("Addon Id:%v Not Publish status", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.Type == consts.PlanTypeRecurringAddon, fmt.Sprintf("Addon Id:%v Not Recurring Type", addon.AddonPlan.Id))
		utility.Assert(addon.AddonPlan.IntervalUnit == plan.IntervalUnit, "update addon must have same recurring interval to plan")
		utility.Assert(addon.AddonPlan.IntervalCount == plan.IntervalCount, "update addon must have same recurring interval to plan")
		TotalAmountExcludingTax = TotalAmountExcludingTax + addon.AddonPlan.Amount*addon.Quantity
	}

	var recurringDiscountCode string
	if len(req.DiscountCode) > 0 {
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:   req.MerchantId,
			UserId:       req.UserId,
			DiscountCode: req.DiscountCode,
			Currency:     plan.Currency,
			PLanId:       plan.Id,
		})
		if canApply {
			if isRecurring {
				recurringDiscountCode = req.DiscountCode
			}
		} else {
			req.DiscountCode = ""
			discountMessage = message
		}
		if req.IsSubmit {
			utility.Assert(canApply, message)
		}
		//if isRecurring {
		//	recurringDiscountCode = req.DiscountCode
		//}
	}

	var currentTimeStart = gtime.Now()
	var trialEnd = currentTimeStart.Timestamp() - 1
	var cancelAtPeriodEnd = 0
	if plan.TrialDurationTime > 0 || req.TrialEnd > 0 {
		var totalAmountExcludingTax = plan.TrialAmount * req.Quantity
		if plan.TrialDurationTime > 0 && req.TrialEnd == 0 {
			req.TrialEnd = currentTimeStart.Timestamp() + plan.TrialDurationTime
		} else {
			// if trialEnd set, ignore plan.TrialAmount and plan.demand
			totalAmountExcludingTax = 0
		}
		//trial period
		if plan.TrialAmount > 0 {
			utility.Assert(len(addons) == 0, "addon is not available for charge trial plan")
		}

		cancelAtPeriodEnd = plan.CancelAtTrialEnd
		var currentTimeEnd = req.TrialEnd
		trialEnd = currentTimeEnd
		discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, plan.MerchantId, totalAmountExcludingTax, plan.Currency, req.DiscountCode, currentTimeStart.Timestamp()), totalAmountExcludingTax)
		totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
		var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(subscriptionTaxPercentage))
		invoice := &bean.InvoiceSimplify{
			InvoiceName:                    "SubscriptionCreate",
			ProductName:                    plan.PlanName,
			OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount,
			TotalAmount:                    totalAmountExcludingTax + taxAmount,
			TotalAmountExcludingTax:        totalAmountExcludingTax,
			DiscountCode:                   req.DiscountCode,
			DiscountAmount:                 discountAmount,
			Currency:                       plan.Currency,
			TaxAmount:                      taxAmount,
			BizType:                        consts.BizTypeSubscription,
			SubscriptionAmount:             totalAmountExcludingTax + discountAmount + taxAmount,
			SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount,
			TrialEnd:                       trialEnd,
			Lines: []*bean.InvoiceItemSimplify{{
				Currency:               plan.Currency,
				OriginAmount:           totalAmountExcludingTax + taxAmount + discountAmount,
				Amount:                 totalAmountExcludingTax + taxAmount,
				DiscountAmount:         discountAmount,
				Tax:                    taxAmount,
				AmountExcludingTax:     totalAmountExcludingTax,
				TaxPercentage:          subscriptionTaxPercentage,
				UnitAmountExcludingTax: plan.TrialAmount,
				Description:            plan.Description,
				Proration:              false,
				Quantity:               req.Quantity,
				PeriodEnd:              currentTimeEnd,
				PeriodStart:            currentTimeStart.Timestamp(),
				Plan:                   bean.SimplifyPlan(plan),
			}},
			PeriodStart: currentTimeStart.Timestamp(),
			PeriodEnd:   currentTimeEnd,
			FinishTime:  currentTimeStart.Timestamp(),
		}
		return &CreatePreviewInternalRes{
			Plan:                     plan,
			TrialEnd:                 trialEnd,
			Quantity:                 req.Quantity,
			Gateway:                  gateway,
			Merchant:                 merchantInfo,
			AddonParams:              req.AddonParams,
			Addons:                   addons,
			OriginAmount:             invoice.OriginAmount,
			TotalAmount:              invoice.TotalAmount,
			DiscountAmount:           invoice.DiscountAmount,
			Invoice:                  invoice,
			RecurringDiscountCode:    recurringDiscountCode,
			Discount:                 bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, invoice.DiscountCode)),
			Currency:                 currency,
			VatCountryCode:           vatCountryCode,
			VatCountryName:           vatCountryName,
			VatNumber:                req.VatNumber,
			VatNumberValidate:        vatNumberValidate,
			VatVerifyData:            utility.MarshalToJsonString(vatNumberValidate),
			UserId:                   req.UserId,
			Email:                    user.Email,
			VatCountryRate:           vatCountryRate,
			Gateways:                 service2.GetMerchantAvailableGatewaysByCountryCode(ctx, req.MerchantId, req.VatCountryCode),
			TaxPercentage:            subscriptionTaxPercentage,
			VatNumberValidateMessage: vatNumberValidateMessage,
			DiscountMessage:          discountMessage,
		}, nil
	} else {
		var currentTimeEnd = subscription2.GetPeriodEndFromStart(ctx, currentTimeStart.Timestamp(), req.PlanId)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			InvoiceName:   "SubscriptionCreate",
			DiscountCode:  req.DiscountCode,
			TimeNow:       gtime.Now().Timestamp(),
			Currency:      currency,
			PlanId:        req.PlanId,
			Quantity:      req.Quantity,
			AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
			TaxPercentage: subscriptionTaxPercentage,
			PeriodStart:   currentTimeStart.Timestamp(),
			PeriodEnd:     currentTimeEnd,
			FinishTime:    currentTimeStart.Timestamp(),
		})

		return &CreatePreviewInternalRes{
			Plan:                     plan,
			TrialEnd:                 trialEnd,
			Quantity:                 req.Quantity,
			Gateway:                  gateway,
			Merchant:                 merchantInfo,
			AddonParams:              req.AddonParams,
			Addons:                   addons,
			OriginAmount:             invoice.OriginAmount,
			TotalAmount:              invoice.TotalAmount,
			DiscountAmount:           invoice.DiscountAmount,
			Invoice:                  invoice,
			RecurringDiscountCode:    recurringDiscountCode,
			Discount:                 bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, invoice.DiscountCode)),
			Currency:                 currency,
			VatCountryCode:           vatCountryCode,
			VatCountryName:           vatCountryName,
			VatNumber:                req.VatNumber,
			VatNumberValidate:        vatNumberValidate,
			VatVerifyData:            utility.MarshalToJsonString(vatNumberValidate),
			UserId:                   req.UserId,
			Email:                    user.Email,
			VatCountryRate:           vatCountryRate,
			Gateways:                 service2.GetMerchantAvailableGatewaysByCountryCode(ctx, req.MerchantId, req.VatCountryCode),
			TaxPercentage:            subscriptionTaxPercentage,
			VatNumberValidateMessage: vatNumberValidateMessage,
			DiscountMessage:          discountMessage,
			CancelAtPeriodEnd:        cancelAtPeriodEnd,
		}, nil
	}
}

func SubscriptionCreate(ctx context.Context, req *CreateInternalReq) (*CreateInternalRes, error) {
	if req.Discount != nil {
		//utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call") // todo mark enable for test automatic
		// create external discount
		utility.Assert(req.PlanId > 0, "PlanId invalid")
		utility.Assert(req.GatewayId > 0, "GatewayId invalid")
		utility.Assert(req.UserId > 0, "UserId invalid")
		plan := query.GetPlanById(ctx, req.PlanId)
		utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
		utility.Assert(plan != nil, "invalid planId")
		one := discount.CreateExternalDiscount(ctx, req.MerchantId, req.UserId, strconv.FormatUint(req.PlanId, 10), req.Discount, plan.Currency)
		req.DiscountCode = one.Code
	} else if len(req.DiscountCode) > 0 {
		one := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
		utility.Assert(one.Type == 0, "invalid code, code is from external")
	}

	prepare, err := SubscriptionCreatePreview(ctx, &CreatePreviewInternalReq{
		MerchantId:     req.MerchantId,
		PlanId:         req.PlanId,
		UserId:         req.UserId,
		DiscountCode:   req.DiscountCode,
		Quantity:       req.Quantity,
		GatewayId:      req.GatewayId,
		AddonParams:    req.AddonParams,
		VatCountryCode: req.VatCountryCode,
		VatNumber:      req.VatNumber,
		TaxPercentage:  req.TaxPercentage,
		IsSubmit:       true,
		TrialEnd:       req.TrialEnd,
	})
	if err != nil {
		return nil, err
	}
	// todo mark countryCode is required or node
	// utility.Assert(len(prepare.VatCountryCode) > 0, "CountryCode Needed")
	if req.ConfirmTotalAmount > 0 {
		utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch preview again")
	}
	if len(req.ConfirmCurrency) > 0 {
		utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch preview again")
	}

	//if prepare.Gateway.GatewayType == consts.GatewayTypeWireTransfer {
	//	utility.Assert(prepare.Invoice.TotalAmount >= prepare.Gateway.MinimumAmount, "Total Amount not reach the wire transfer's minimum amount")
	//}

	if prepare.Invoice.TotalAmount == 0 && strings.Contains(prepare.Plan.TrialDemand, "paymentMethod") && req.PaymentMethodId == "" {
		utility.Assert(prepare.Gateway.GatewayType == consts.GatewayTypeDefault, "card payment gateway need") // todo mark
	}

	var subType = consts.SubTypeDefault
	if consts.SubscriptionCycleUnderUniBeeControl {
		subType = consts.SubTypeUniBeeControl
	}

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, prepare.Invoice.PeriodEnd, prepare.Plan.Id)

	one := &entity.Subscription{
		MerchantId:                  prepare.Merchant.Id,
		Type:                        subType,
		PlanId:                      prepare.Plan.Id,
		TrialEnd:                    prepare.TrialEnd,
		GatewayId:                   prepare.Gateway.Id,
		UserId:                      prepare.UserId,
		Quantity:                    prepare.Quantity,
		Amount:                      prepare.TotalAmount, // todo mark should use originAmount
		Currency:                    prepare.Currency,
		AddonData:                   utility.MarshalToJsonString(prepare.AddonParams),
		SubscriptionId:              utility.CreateSubscriptionId(),
		Status:                      consts.SubStatusInit,
		CustomerEmail:               prepare.Email,
		ReturnUrl:                   req.ReturnUrl,
		VatNumber:                   prepare.VatNumber,
		VatVerifyData:               prepare.VatVerifyData,
		CountryCode:                 prepare.VatCountryCode,
		TaxPercentage:               prepare.TaxPercentage,
		CurrentPeriodStart:          prepare.Invoice.PeriodStart,
		CurrentPeriodEnd:            prepare.Invoice.PeriodEnd,
		DunningTime:                 dunningTime,
		BillingCycleAnchor:          prepare.Invoice.PeriodStart,
		GatewayDefaultPaymentMethod: req.PaymentMethodId,
		DiscountCode:                prepare.RecurringDiscountCode,
		CreateTime:                  gtime.Now().Timestamp(),
		MetaData:                    utility.MarshalToJsonString(req.Metadata),
		GasPayer:                    prepare.Plan.GasPayer,
		CancelAtPeriodEnd:           prepare.CancelAtPeriodEnd,
	}

	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(uint(id))

	var createRes *gateway_bean.GatewayCreateSubscriptionResp
	invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, one)
	utility.AssertError(err, "System Error")
	if prepare.Invoice.TotalAmount == 0 {
		//totalAmount is 0, no payment need
		utility.AssertError(err, "System Error")
		if strings.Contains(prepare.Plan.TrialDemand, "paymentMethod") && req.PaymentMethodId == "" {
			url, _ := method.NewPaymentMethod(ctx, &method.NewPaymentMethodInternalReq{
				MerchantId:     _interface.GetMerchantId(ctx),
				UserId:         req.UserId,
				Currency:       prepare.Currency,
				GatewayId:      req.GatewayId,
				SubscriptionId: one.SubscriptionId,
				RedirectUrl:    req.ReturnUrl,
				Metadata:       map[string]interface{}{"InvoiceId": invoice.InvoiceId, "Action": "SubscriptionCreate"},
			})
			createRes = &gateway_bean.GatewayCreateSubscriptionResp{
				GatewaySubscriptionId: one.SubscriptionId,
				Link:                  url,
				Paid:                  false,
			}
		} else {
			invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
			utility.AssertError(err, "System Error")
			createRes = &gateway_bean.GatewayCreateSubscriptionResp{
				GatewaySubscriptionId: one.SubscriptionId,
				Paid:                  true,
			}
		}
		// todo mark subscription become active with payment mq message
	} else if len(req.PaymentMethodId) > 0 {
		// createAndPayNewProrationInvoice
		merchant := query.GetMerchantById(ctx, one.MerchantId)
		utility.Assert(merchant != nil, "merchant not found")
		//utility.Assert(user != nil, "user not found")
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		utility.AssertError(err, "System Error")
		var createPaymentResult, err = service.CreateSubInvoicePaymentDefaultAutomatic(ctx, one.GatewayDefaultPaymentMethod, invoice, gateway.Id, false, req.ReturnUrl, "SubscriptionCreate")
		if err != nil {
			g.Log().Print(ctx, "SubscriptionCreate CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return nil, err
		}
		createRes = &gateway_bean.GatewayCreateSubscriptionResp{
			GatewaySubscriptionId: createPaymentResult.Payment.PaymentId,
			Data:                  utility.MarshalToJsonString(createPaymentResult),
			Link:                  createPaymentResult.Link,
			Paid:                  createPaymentResult.Status == consts.PaymentSuccess,
		}
	} else {
		gateway := query.GetGatewayById(ctx, one.GatewayId)
		if gateway == nil {
			return nil, gerror.New("SubscriptionBillingCycleDunningInvoice gateway not found")
		}
		utility.AssertError(err, "System Error")
		//var createPaymentResult, err = service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
		//	CheckoutMode: true,
		//	Gateway:      prepare.Gateway,
		//	Pay: &entity.Payment{
		//		SubscriptionId:    one.SubscriptionId,
		//		ExternalPaymentId: one.SubscriptionId,
		//		BizType:           consts.BizTypeSubscription,
		//		UserId:            prepare.UserId,
		//		GatewayId:         prepare.Gateway.Id,
		//		TotalAmount:       prepare.Invoice.TotalAmount,
		//		Currency:          prepare.Invoice.Currency,
		//		CryptoAmount:      prepare.Invoice.CryptoAmount,
		//		CryptoCurrency:    prepare.Invoice.CryptoCurrency,
		//		CountryCode:       prepare.VatCountryCode,
		//		MerchantId:        prepare.Merchant.Id,
		//		CompanyId:         prepare.Merchant.CompanyId,
		//		BillingReason:     prepare.Invoice.InvoiceName,
		//		ReturnUrl:         req.ReturnUrl,
		//		GasPayer:          prepare.Plan.GasPayer,
		//	},
		//	ExternalUserId: strconv.FormatUint(one.UserId, 10),
		//	Email:          prepare.Email,
		//	Invoice:        bean.SimplifyInvoice(invoice),
		//	Metadata:       map[string]interface{}{"BillingReason": prepare.Invoice.InvoiceName},
		//})
		createPaymentResult, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, "", invoice, one.GatewayId, true, req.ReturnUrl, "SubscriptionCreate")
		if err != nil {
			_, updateErr := dao.Subscription.Ctx(ctx).Data(g.Map{
				dao.Subscription.Columns().Status:    consts.SubStatusCancelled,
				dao.Subscription.Columns().GmtModify: gtime.Now(),
			}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
			if updateErr != nil {
				return nil, updateErr
			}
			utility.AssertError(err, "Create Payment Error")
		}
		createRes = &gateway_bean.GatewayCreateSubscriptionResp{
			GatewaySubscriptionId: createPaymentResult.Payment.PaymentId,
			Data:                  utility.MarshalToJsonString(createPaymentResult),
			Link:                  createPaymentResult.Link,
			Paid:                  createPaymentResult.Status == consts.PaymentSuccess,
		}
	}

	//Update Subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().GatewaySubscriptionId: createRes.GatewaySubscriptionId,
		dao.Subscription.Columns().Status:                consts.SubStatusPending,
		dao.Subscription.Columns().Link:                  createRes.Link,
		dao.Subscription.Columns().ResponseData:          createRes.Data,
		dao.Subscription.Columns().GmtModify:             gtime.Now(),
	}).Where(dao.Subscription.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	one.GatewaySubscriptionId = createRes.GatewaySubscriptionId
	one.Status = consts.SubStatusPending
	one.Link = createRes.Link

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCreate.Topic,
		Tag:   redismq2.TopicSubscriptionCreate.Tag,
		Body:  one.SubscriptionId,
	})
	if createRes.Paid {
		utility.Assert(invoice.Id > 0, "Server Error")
		oneInvoice := query.GetInvoiceByInvoiceId(ctx, invoice.InvoiceId)
		err = handler.HandleSubscriptionFirstInvoicePaid(ctx, one, oneInvoice)
		one = query.GetSubscriptionBySubscriptionId(ctx, one.SubscriptionId)
		utility.AssertError(err, "Finish Subscription Error")
	} else if req.StartIncomplete {
		err = SubscriptionActiveTemporarily(ctx, one.SubscriptionId, one.CurrentPeriodEnd)
		utility.AssertError(err, "Start Active Temporarily")
	}
	return &CreateInternalRes{
		Subscription: bean.SimplifySubscription(one),
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}

type UpdatePreviewInternalRes struct {
	Subscription          *entity.Subscription               `json:"subscription"`
	Plan                  *entity.Plan                       `json:"plan"`
	Quantity              int64                              `json:"quantity"`
	Gateway               *entity.MerchantGateway            `json:"gateway"`
	MerchantInfo          *entity.Merchant                   `json:"merchantInfo"`
	AddonParams           []*bean.PlanAddonParam             `json:"addonParams"`
	Addons                []*bean.PlanAddonDetail            `json:"addons"`
	OriginAmount          int64                              `json:"originAmount"                `
	TotalAmount           int64                              `json:"totalAmount"`
	DiscountAmount        int64                              `json:"discountAmount"`
	Currency              string                             `json:"currency"`
	UserId                uint64                             `json:"userId"`
	OldPlan               *entity.Plan                       `json:"oldPlan"`
	Invoice               *bean.InvoiceSimplify              `json:"invoice"`
	NextPeriodInvoice     *bean.InvoiceSimplify              `json:"nextPeriodInvoice"`
	ProrationDate         int64                              `json:"prorationDate"`
	EffectImmediate       bool                               `json:"EffectImmediate"`
	Gateways              []*bean.GatewaySimplify            `json:"gateways"`
	Changed               bool                               `json:"changed"`
	IsUpgrade             bool                               `json:"isUpgrade"`
	TaxPercentage         int64                              `json:"taxPercentage" `
	RecurringDiscountCode string                             `json:"recurringDiscountCode" `
	Discount              *bean.MerchantDiscountCodeSimplify `json:"discount" `
}

func isUpgradeForSubscription(ctx context.Context, sub *entity.Subscription, plan *entity.Plan, quantity int64, addonParams []*bean.PlanAddonParam) (isUpgrade bool, changed bool) {
	//default logical，Effect Immediately for upgrade, effect at period end for downgrade
	//situation 1，NewPlan Unit Amount >  OldPlan Unit Amount，is upgrade，ignore Quantity and addon change
	//situation 2，NewPlan Unit Amount <  OldPlan Unit Amount，is downgrade，ignore Quantity and addon change
	//situation 3，NewPlan Total Amount >  OldPlan Total Amount，is upgrade
	//situation 4，NewPlan Total Amount <  OldPlan Total Amount，is downgrade
	//situation 5，NewPlan Total Amount =  OldPlan Total Amount，see Addon changes，if new addon appended or addon quantity changed, is upgrade，otherwise downgrade
	oldPlan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(oldPlan != nil, "oldPlan not found")
	if plan.IntervalUnit != oldPlan.IntervalUnit || plan.IntervalCount != oldPlan.IntervalCount {
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
		utility.Assert(err == nil, fmt.Sprintf("UnmarshalFromJsonString internal err:%v", err))
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
	GatewayId       uint64                 `json:"gatewayId" dc:"Id" `
	EffectImmediate int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
	TaxPercentage   *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
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
	var gatewayId = sub.GatewayId
	if req.GatewayId > 0 {
		gatewayId = req.GatewayId
	}
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(service2.IsGatewaySupportCountryCode(ctx, gateway, sub.CountryCode), "gateway not support")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	utility.Assert(sub.CancelAtPeriodEnd == 0, "subscription will cancel at period end, should resume subscription first")
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	addons := checkAndListAddonsFromParams(ctx, req.AddonParams)
	var subscriptionTaxPercentage = sub.TaxPercentage
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
	if isUpgrade {
		effectImmediate = true
	} else {
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

	var currentInvoice *bean.InvoiceSimplify
	var nextPeriodInvoice *bean.InvoiceSimplify
	var RecurringDiscountCode string
	if prorationDate == 0 {
		prorationDate = time.Now().Unix()
		if sub.TestClock > sub.CurrentPeriodStart && !config2.GetConfigInstance().IsProd() {
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
		})
		utility.Assert(canApply, message)
		if isRecurring {
			RecurringDiscountCode = req.DiscountCode
		}
	} else if len(sub.DiscountCode) > 0 {
		canApply, isRecurring, _ := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:     sub.MerchantId,
			UserId:         sub.UserId,
			DiscountCode:   sub.DiscountCode,
			Currency:       sub.Currency,
			SubscriptionId: sub.SubscriptionId,
			PLanId:         req.NewPlanId,
		})
		if canApply && isRecurring {
			req.DiscountCode = sub.DiscountCode
			RecurringDiscountCode = sub.DiscountCode
		}
	}

	if effectImmediate {
		if !config.GetMerchantSubscriptionConfig(ctx, sub.MerchantId).UpgradeProration {
			// without proration, just generate next cycle
			currentInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:   "SubscriptionCycle",
				Currency:      sub.Currency,
				DiscountCode:  req.DiscountCode,
				TimeNow:       prorationDate,
				PlanId:        req.NewPlanId,
				Quantity:      req.Quantity,
				AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
				TaxPercentage: subscriptionTaxPercentage,
				PeriodStart:   prorationDate,
				PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, prorationDate, req.NewPlanId),
				FinishTime:    prorationDate,
			})
		} else if prorationDate < sub.CurrentPeriodStart {
			// after period end before trial end, also or sub data not sync or use testClock in stage env
			currentInvoice = &bean.InvoiceSimplify{
				InvoiceName:                    "SubscriptionUpgrade",
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
				FinishTime:                     prorationDate,
			}
		} else if prorationDate > sub.CurrentPeriodEnd {
			// after periodEnd, is not a currentInvoice, just use it
			currentInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:   "SubscriptionCycle",
				Currency:      sub.Currency,
				DiscountCode:  req.DiscountCode,
				TimeNow:       prorationDate,
				PlanId:        req.NewPlanId,
				Quantity:      req.Quantity,
				AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
				TaxPercentage: subscriptionTaxPercentage,
				PeriodStart:   prorationDate,
				PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, prorationDate, req.NewPlanId),
				FinishTime:    prorationDate,
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
			if !hasIntervalChange {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToFixedEndInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:       "SubscriptionUpgrade",
					ProductName:       plan.PlanName,
					Currency:          sub.Currency,
					DiscountCode:      req.DiscountCode,
					TimeNow:           prorationDate,
					TaxPercentage:     subscriptionTaxPercentage,
					ProrationDate:     prorationDate,
					OldProrationPlans: oldProrationPlanParams,
					NewProrationPlans: newProrationPlanParams,
					PeriodStart:       sub.CurrentPeriodStart,
					PeriodEnd:         sub.CurrentPeriodEnd,
				})
			} else {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToDifferentIntervalInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:       "SubscriptionUpgrade",
					ProductName:       plan.PlanName,
					Currency:          sub.Currency,
					DiscountCode:      req.DiscountCode,
					TimeNow:           prorationDate,
					TaxPercentage:     subscriptionTaxPercentage,
					ProrationDate:     prorationDate,
					OldProrationPlans: oldProrationPlanParams,
					NewProrationPlans: newProrationPlanParams,
					PeriodStart:       sub.CurrentPeriodStart,
					PeriodEnd:         sub.CurrentPeriodEnd,
				})
			}
		}
		prorationDate = currentInvoice.ProrationDate
	} else {
		prorationDate = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
		currentInvoice = &bean.InvoiceSimplify{
			InvoiceName:                    "SubscriptionUpgrade",
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
		}
	}

	nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:   "SubscriptionCycle",
		Currency:      sub.Currency,
		DiscountCode:  req.DiscountCode,
		TimeNow:       prorationDate,
		PlanId:        req.NewPlanId,
		Quantity:      req.Quantity,
		AddonJsonData: utility.MarshalToJsonString(req.AddonParams),
		TaxPercentage: subscriptionTaxPercentage,
		PeriodStart:   utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd),
		PeriodEnd:     subscription2.GetPeriodEndFromStart(ctx, utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd), req.NewPlanId),
		FinishTime:    utility.MaxInt64(currentInvoice.PeriodEnd, sub.TrialEnd),
	})

	if currentInvoice.TotalAmount <= 0 {
		effectImmediate = false
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
		RecurringDiscountCode: RecurringDiscountCode,
		Discount:              bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, currentInvoice.DiscountCode)),
	}, nil
}

type UpdateSubscriptionInternalResp struct {
	GatewayUpdateId string          `json:"gatewayUpdateId" description:""`
	Data            string          `json:"data"`
	Link            string          `json:"link" description:""`
	Paid            bool            `json:"paid" description:""`
	Invoice         *entity.Invoice `json:"invoice" description:""`
}

type UpdateInternalReq struct {
	SubscriptionId     string                      `json:"subscriptionId" dc:"SubscriptionId" v:"required"`
	NewPlanId          uint64                      `json:"newPlanId" dc:"NewPlanId" v:"required"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          uint64                      `json:"gatewayId" dc:"GatewayId" `
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
}

func SubscriptionUpdate(ctx context.Context, req *UpdateInternalReq, merchantMemberId int64) (*subscription.UpdateRes, error) {
	prorationDate := gtime.Now().Timestamp()
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
		one := discount.CreateExternalDiscount(ctx, sub.MerchantId, sub.UserId, strconv.FormatUint(req.NewPlanId, 10), req.Discount, plan.Currency)
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
	}, prorationDate, merchantMemberId)
	if err != nil {
		return nil, err
	}

	//subscription prepare
	if req.ConfirmTotalAmount > 0 {
		utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch preview again")
	}
	if len(req.ConfirmCurrency) > 0 {
		utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch again")
	}
	if prepare.Invoice.TotalAmount <= 0 {
		utility.Assert(prepare.EffectImmediate == false, "System Error, Cannot Effect Immediate With Negative CaptureAmount")
	}

	var effectImmediate = 0
	var effectTime = prepare.Subscription.CurrentPeriodEnd
	if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
		effectImmediate = 1
		effectTime = gtime.Now().Timestamp()
	}

	gatewayId := prepare.Subscription.GatewayId
	if req.GatewayId > 0 {
		gatewayId = req.GatewayId
	}

	one := &entity.SubscriptionPendingUpdate{
		MerchantId:       prepare.MerchantInfo.Id,
		GatewayId:        gatewayId,
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
		MetaData:         utility.MarshalToJsonString(req.Metadata),
	}

	result, err := dao.SubscriptionPendingUpdate.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)
	var subUpdateRes *UpdateSubscriptionInternalResp
	if prepare.EffectImmediate && prepare.Invoice.TotalAmount > 0 {
		// createAndPayNewProrationInvoice
		merchantInfo := query.GetMerchantById(ctx, one.MerchantId)
		utility.Assert(merchantInfo != nil, "merchantInfo not found")
		// utility.Assert(user != nil, "user not found")
		gateway := query.GetGatewayById(ctx, gatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, prepare.Subscription)
		utility.AssertError(err, "System Error")
		createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, prepare.Subscription.GatewayDefaultPaymentMethod, invoice, gateway.Id, req.ManualPayment, req.ReturnUrl, "SubscriptionUpdate")
		if err != nil {
			g.Log().Print(ctx, "SubscriptionUpdate CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return nil, err
		}
		// Upgrade
		subUpdateRes = &UpdateSubscriptionInternalResp{
			GatewayUpdateId: createRes.Invoice.InvoiceId,
			Data:            utility.MarshalToJsonString(createRes),
			Link:            createRes.Link,
			Paid:            createRes.Status == consts.PaymentSuccess,
			Invoice:         createRes.Invoice,
		}
	} else if prepare.EffectImmediate && prepare.Invoice.TotalAmount == 0 {
		//totalAmount is 0, no payment need
		invoice, err := handler2.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, prepare.Subscription)
		utility.AssertError(err, "System Error")
		invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
		utility.AssertError(err, "System Error")
		subUpdateRes = &UpdateSubscriptionInternalResp{
			Paid: true,
			Link: "",
		}
	} else {
		prepare.EffectImmediate = false
		subUpdateRes = &UpdateSubscriptionInternalResp{
			Paid: false,
			Link: "",
		}
	}

	if err != nil {
		return nil, err
	}
	// bing to subscription
	_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().DiscountCode:    prepare.RecurringDiscountCode,
		dao.Subscription.Columns().PendingUpdateId: one.PendingUpdateId,
		dao.Subscription.Columns().GmtModify:       gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, one.SubscriptionId).OmitNil().Update()
	if err != nil {
		return nil, err
	}
	// only one need, cancel others
	// need cancel payment、 invoice and send invoice email
	CancelOtherUnfinishedPendingUpdatesBackground(prepare.Subscription.SubscriptionId, one.PendingUpdateId, "CancelByNewUpdate-"+one.PendingUpdateId)

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
		dao.SubscriptionPendingUpdate.Columns().Status:       consts.PendingSubStatusCreate,
		dao.SubscriptionPendingUpdate.Columns().ResponseData: subUpdateRes.Data,
		dao.SubscriptionPendingUpdate.Columns().GmtModify:    gtime.Now(),
		dao.SubscriptionPendingUpdate.Columns().Paid:         PaidInt,
		dao.SubscriptionPendingUpdate.Columns().Link:         subUpdateRes.Link,
		dao.SubscriptionPendingUpdate.Columns().InvoiceId:    subUpdateRes.GatewayUpdateId,
		dao.SubscriptionPendingUpdate.Columns().Note:         note,
		dao.SubscriptionPendingUpdate.Columns().MetaData:     utility.MarshalToJsonString(req.Metadata),
	}).Where(dao.SubscriptionPendingUpdate.Columns().PendingUpdateId, one.PendingUpdateId).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	if prepare.EffectImmediate && subUpdateRes.Paid {
		_, err = handler.HandlePendingUpdatePaymentSuccess(ctx, prepare.Subscription, one.PendingUpdateId, subUpdateRes.Invoice)
		if err != nil {
			return nil, err
		}
		one.Status = consts.PendingSubStatusFinished
	}

	return &subscription.UpdateRes{
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
			GmtModify:       one.GmtModify,
			Paid:            one.Paid,
			Link:            one.Link,
			MerchantMember:  bean.SimplifyMerchantMember(query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
			EffectImmediate: one.EffectImmediate,
			EffectTime:      one.EffectTime,
			Note:            one.Note,
			Plan:            bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:          addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UpdatePlan:      bean.SimplifyPlan(query.GetPlanById(ctx, one.UpdatePlanId)),
			UpdateAddons:    addon2.GetSubscriptionAddonsByAddonJson(ctx, one.UpdateAddonData),
			Metadata:        req.Metadata,
		},
		Paid: len(subUpdateRes.Link) == 0 || subUpdateRes.Paid, // link is blank or paid is true, portal will not redirect
		Link: subUpdateRes.Link,
		Note: note,
	}, nil
}

func SubscriptionCancel(ctx context.Context, subscriptionId string, proration bool, invoiceNow bool, reason string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	if sub.Status == consts.SubStatusCancelled || sub.Status == consts.SubStatusExpired {
		g.Log().Infof(ctx, "SubscriptionCancel, subscription already cancelled or expired")
		return nil
	}
	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	if !config2.GetConfigInstance().IsServerDev() || !config2.GetConfigInstance().IsLocal() {
		// todo mark will support proration invoiceNow later
		invoiceNow = false
		proration = false
		// todo mark will support proration invoiceNow later
		// only local env can cancel immediately invoice_compute proration invoice
		utility.Assert(invoiceNow == false && proration == false, "cancel subscription with proration invoice immediate not support for this version")
	}
	var nextStatus = consts.SubStatusCancelled
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:       nextStatus,
		dao.Subscription.Columns().CancelReason: reason,
		dao.Subscription.Columns().TrialEnd:     sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:    gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, sub.UserId)
	if user != nil {
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
		if merchant != nil {
			var template = email.TemplateSubscriptionImmediateCancel
			if (sub.Status == consts.SubStatusIncomplete || sub.Status == consts.SubStatusActive) && sub.TrialEnd >= sub.CurrentPeriodEnd {
				//first trial period without payment
				template = email.TemplateSubscriptionCancelledByTrialEnd
			}
			err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, template, "", &email.TemplateVariable{
				UserName:            user.FirstName + " " + user.LastName,
				MerchantProductName: plan.PlanName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
				PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancel:%s", err.Error())
			}
		}
	}

	_, _ = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicSubscriptionCancel.Topic,
		Tag:   redismq2.TopicSubscriptionCancel.Tag,
		Body:  sub.SubscriptionId,
	})
	return nil
}

func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool, merchantMemberId int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 1 {
		//已经设置未周期结束取消
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 1,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}

	user := query.GetUserAccountById(ctx, sub.UserId)
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	// SendEmail
	if merchantMemberId > 0 {
		//merchant Cancel
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelledAtPeriodEndByMerchantAdmin, "", &email.TemplateVariable{
			UserName:            user.FirstName + " " + user.LastName,
			MerchantProductName: plan.PlanName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
		})
		if err != nil {
			g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelAtPeriodEnd:%s", err.Error())
		}
	} else {
		err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelledAtPeriodEndByUser, "", &email.TemplateVariable{
			UserName:            user.FirstName + " " + user.LastName,
			MerchantProductName: plan.PlanName,
			MerchantCustomEmail: merchant.Email,
			MerchantName:        query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
			PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
		})
		if err != nil {
			g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelAtPeriodEnd:%s", err.Error())
		}
	}
	return nil
}

func SubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 0 {
		//已经设置未周期结束取消
		return nil
	}

	plan := query.GetPlanById(ctx, sub.PlanId)
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CancelAtPeriodEnd: 0,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	user := query.GetUserAccountById(ctx, sub.UserId)
	merchant := query.GetMerchantById(ctx, sub.MerchantId)
	err = email.SendTemplateEmail(ctx, merchant.Id, user.Email, user.TimeZone, email.TemplateSubscriptionCancelLastCancelledAtPeriodEnd, "", &email.TemplateVariable{
		UserName:            user.FirstName + " " + user.LastName,
		MerchantProductName: plan.PlanName,
		MerchantCustomEmail: merchant.Email,
		MerchantName:        query.GetMerchantCountryConfigName(ctx, merchant.Id, user.CountryCode),
		PeriodEnd:           gtime.NewFromTimeStamp(sub.CurrentPeriodEnd),
	})
	if err != nil {
		g.Log().Errorf(ctx, "SendTemplateEmail SubscriptionCancelLastCancelAtPeriodEnd:%s", err.Error())
	}
	return nil
}

//
//type AdminAttachSubscriptionToUserEmailReq struct {
//	ExternalUserId string                      `json:"externalUserId" dc:"ExternalUserId"`
//	Email          string                      `json:"email" dc:"Email" v:"required"`
//	MerchantId     uint64                      `json:"merchantId" dc:"MerchantId" v:"required"`
//	PlanId         uint64                      `json:"planId" dc:"PlanId" v:"required"`
//	UserId         uint64                      `json:"userId" dc:"UserId" v:"required"`
//	Quantity       int64                       `json:"quantity" dc:"Quantity，Default 1" `
//	GatewayId      uint64                      `json:"gatewayId" dc:"Id"   v:"required" `
//	AddonParams    []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
//}
//
//func AdminAttachSubscriptionToUserEmail(ctx context.Context, req *AdminAttachSubscriptionToUserEmailReq) (*entity.Subscription, error) {
//	user, err := auth.QueryOrCreateUser(ctx, &auth.NewReq{
//		ExternalUserId: req.ExternalUserId,
//		Email:          req.Email,
//		MerchantId:     req.MerchantId,
//	})
//	utility.AssertError(err, "QueryOrCreateUser error")
//	var subType = consts.SubTypeDefault
//	if consts.SubscriptionCycleUnderUniBeeControl {
//		subType = consts.SubTypeUniBeeControl
//	}
//	one := &entity.Subscription{
//		MerchantId:                  req.MerchantId,
//		Type:                        subType,
//		PlanId:                      prepare.Plan.Id,
//		TrialEnd:                    prepare.TrialEnd,
//		GatewayId:                   prepare.Gateway.Id,
//		UserId:                      prepare.UserId,
//		Quantity:                    prepare.Quantity,
//		Amount:                      prepare.TotalAmount,
//		Currency:                    prepare.Currency,
//		AddonData:                   utility.MarshalToJsonString(prepare.AddonParams),
//		SubscriptionId:              utility.CreateSubscriptionId(),
//		Status:                      consts.SubStatusPending,
//		CustomerEmail:               prepare.Email,
//		ReturnUrl:                   req.ReturnUrl,
//		VatNumber:                   prepare.VatNumber,
//		VatVerifyData:               prepare.VatVerifyData,
//		CountryCode:                 prepare.VatCountryCode,
//		TaxPercentage:               prepare.TaxPercentage,
//		CurrentPeriodStart:          prepare.Invoice.PeriodStart,
//		CurrentPeriodEnd:            prepare.Invoice.PeriodEnd,
//		DunningTime:                 dunningTime,
//		BillingCycleAnchor:          prepare.Invoice.PeriodStart,
//		GatewayDefaultPaymentMethod: req.PaymentMethodId,
//		DiscountCode:                prepare.RecurringDiscountCode,
//		CreateTime:                  gtime.Now().Timestamp(),
//		MetaData:                    utility.MarshalToJsonString(req.Metadata),
//		GasPayer:                    prepare.Plan.GasPayer,
//	}
//
//	result, err := dao.Subscription.Ctx(ctx).Data(one).OmitNil().Insert(one)
//	if err != nil {
//		err = gerror.Newf(`SubscriptionCreate record insert failure %s`, err)
//		return nil, err
//	}
//	id, _ := result.LastInsertId()
//	one.Id = uint64(uint(id))
//	return nil, gerror.New("not support")
//}

func SubscriptionAddNewTrialEnd(ctx context.Context, subscriptionId string, AppendNewTrialEndByHour int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	//utility.Assert(sub.Status != consts.SubStatusExpired && sub.Status != consts.SubStatusCancelled, "sub cancelled or sub expired")
	//utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")

	utility.Assert(AppendNewTrialEndByHour > 0, "invalid AppendNewTrialEndByHour , should > 0")
	newTrialEnd := sub.CurrentPeriodEnd + AppendNewTrialEndByHour*3600

	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, sub.PlanId)
	newStatus := sub.Status
	if newTrialEnd > gtime.Now().Timestamp() {
		//automatic change sub status to active
		newStatus = consts.SubStatusActive
		if sub.Status != consts.SubStatusActive {
			service3.TryCancelSubscriptionLatestInvoice(ctx, sub)
		}
	}
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:             newStatus,
		dao.Subscription.Columns().TrialEnd:           newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
		dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
		dao.Subscription.Columns().GmtModify:          gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	if sub.Status != consts.SubStatusActive {
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicSubscriptionActiveWithoutPayment.Topic,
			Tag:   redismq2.TopicSubscriptionActiveWithoutPayment.Tag,
			Body:  sub.SubscriptionId,
		})
	}
	return nil
}

func SubscriptionActiveTemporarily(ctx context.Context, subscriptionId string, expireTime int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing, "subscription not in pending or processing status")
	utility.Assert(sub.CurrentPeriodStart < expireTime, "expireTime should greater then subscription's period start time")
	utility.Assert(sub.CurrentPeriodEnd >= expireTime, "expireTime should lower then subscription's period end time")
	err := handler.MakeSubscriptionIncomplete(ctx, subscriptionId)
	if err != nil {
		return err
	}

	if sub.TrialEnd > 0 && sub.TrialEnd > sub.CurrentPeriodStart {
		// trial start
		oneUser := query.GetUserAccountById(ctx, sub.UserId)
		plan := query.GetPlanById(ctx, sub.PlanId)
		merchant := query.GetMerchantById(ctx, sub.MerchantId)
		if oneUser != nil && plan != nil && merchant != nil {
			err := email.SendTemplateEmail(ctx, sub.MerchantId, oneUser.Email, oneUser.TimeZone, email.TemplateSubscriptionTrialStart, "", &email.TemplateVariable{
				UserName:            oneUser.FirstName + " " + oneUser.LastName,
				MerchantProductName: plan.PlanName,
				MerchantCustomEmail: merchant.Email,
				MerchantName:        query.GetMerchantCountryConfigName(ctx, sub.MerchantId, oneUser.CountryCode),
			})
			if err != nil {
				g.Log().Errorf(ctx, "SendTemplateEmail TemplateSubscriptionTrialStart:%s", err.Error())
			}
		}
	}

	return nil
}

func SubscriptionEndTrial(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	plan := query.GetPlanById(ctx, sub.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	gateway := query.GetGatewayById(ctx, sub.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not trialed")
	err := EndTrialManual(ctx, sub.SubscriptionId)
	if err != nil {
		return err
	}

	return nil
}

func EndTrialManual(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not in trial period")
	newTrialEnd := sub.CurrentPeriodStart - 1
	var newBillingCycleAnchor = utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd)
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, newBillingCycleAnchor, sub.PlanId)
	newStatus := sub.Status
	if gtime.Now().Timestamp() > sub.CurrentPeriodEnd {
		// todo mark has unfinished pending update
		newStatus = consts.SubStatusIncomplete
		// Payment Pending Enter Incomplete
		plan := query.GetPlanById(ctx, sub.PlanId)

		var nextPeriodStart = gtime.Now().Timestamp()
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, plan.Id)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			Currency:      sub.Currency,
			PlanId:        sub.PlanId,
			Quantity:      sub.Quantity,
			AddonJsonData: sub.AddonData,
			TaxPercentage: sub.TaxPercentage,
			PeriodStart:   nextPeriodStart,
			PeriodEnd:     nextPeriodEnd,
			InvoiceName:   "SubscriptionCycle",
			FinishTime:    nextPeriodStart,
		})
		one, err := handler2.CreateProcessingInvoiceForSub(ctx, invoice, sub)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateProcessingInvoiceForSub err:", err.Error())
			return err
		}
		createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, sub.GatewayDefaultPaymentMethod, one, sub.GatewayId, false, "", "SubscriptionEndTrialManual")
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return err
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().CurrentPeriodStart: invoice.PeriodStart,
			dao.Subscription.Columns().CurrentPeriodEnd:   invoice.PeriodEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
		g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic:", utility.MarshalToJsonString(createRes))
		err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, gtime.Now().Timestamp())
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual HandleSubscriptionIncomplete err:", err.Error())
			return err
		}
	} else {
		_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Status:             newStatus,
			dao.Subscription.Columns().TrialEnd:           newTrialEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().BillingCycleAnchor: newBillingCycleAnchor,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
	}
	return nil
}

func MarkSubscriptionProcessed(ctx context.Context, subscriptionId string) error {
	utility.Assert(len(subscriptionId) > 0, "invalid subscriptionId")
	one := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(one != nil, "subscription not found")
	utility.Assert(one.Status == consts.SubStatusPending, "sub not pending status")
	gateway := query.GetGatewayById(ctx, one.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(gateway.GatewayType == consts.GatewayTypeWireTransfer, "not wire transfer type of subscription")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().Status:    consts.SubStatusProcessing,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}
