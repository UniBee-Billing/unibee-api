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
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/cmd/i18n"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/credit/config"
	"unibee/internal/logic/credit/payment"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/gateway/gateway_bean"
	service2 "unibee/internal/logic/gateway/service"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/method"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/subscription/timeline"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type CreatePreviewInternalReq struct {
	MerchantId             uint64                 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId                 uint64                 `json:"planId" dc:"PlanId" v:"required"`
	UserId                 uint64                 `json:"userId" dc:"UserId"`
	Quantity               int64                  `json:"quantity" dc:"Quantity" `
	DiscountCode           string                 `json:"discountCode"        dc:"DiscountCode"`
	GatewayId              *uint64                `json:"gatewayId" dc:"Id"`
	GatewayPaymentType     string                 `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	AddonParams            []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode         string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber              string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage          *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	TrialEnd               int64                  `json:"trialEnd"  description:"trial_end, utc time"` // trial_end, utc time
	IsSubmit               bool
	ProductData            *bean.PlanProductParam `json:"productData"  dc:"ProductData"  `
	PaymentMethodId        string
	Metadata               map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	ApplyPromoCredit       *bool                  `json:"applyPromoCredit" `
	ApplyPromoCreditAmount *int64                 `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type CreatePreviewInternalRes struct {
	Plan                      *entity.Plan               `json:"plan"`
	User                      *bean.UserAccount          `json:"user"`
	Quantity                  int64                      `json:"quantity"`
	Gateway                   *entity.MerchantGateway    `json:"gateway"`
	Merchant                  *entity.Merchant           `json:"merchantInfo"`
	AddonParams               []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                    []*bean.PlanAddonDetail    `json:"addons"`
	OriginAmount              int64                      `json:"originAmount"                `
	TotalAmount               int64                      `json:"totalAmount" `
	DiscountAmount            int64                      `json:"discountAmount"`
	Currency                  string                     `json:"currency" `
	VatCountryCode            string                     `json:"vatCountryCode" `
	VatCountryName            string                     `json:"vatCountryName" `
	VatNumber                 string                     `json:"vatNumber" `
	VatNumberValidate         *bean.ValidResult          `json:"vatNumberValidate" `
	TaxPercentage             int64                      `json:"taxPercentage" `
	TrialEnd                  int64                      `json:"trialEnd" `
	VatVerifyData             string                     `json:"vatVerifyData" `
	Invoice                   *bean.Invoice              `json:"invoice"`
	UserId                    uint64                     `json:"userId" `
	Email                     string                     `json:"email" `
	VatCountryRate            *bean.VatCountryRate       `json:"vatCountryRate" `
	Gateways                  []*detail.Gateway          `json:"gateways" `
	RecurringDiscountCode     string                     `json:"recurringDiscountCode" `
	Discount                  *bean.MerchantDiscountCode `json:"discount" `
	VatNumberValidateMessage  string                     `json:"vatNumberValidateMessage" `
	DiscountMessage           string                     `json:"discountMessage" `
	CancelAtPeriodEnd         int                        `json:"cancelAtPeriodEnd"           description:"whether cancel at period end，0-false | 1-true"` // whether cancel at period end，0-false | 1-true
	GatewayPaymentMethodId    string
	GatewayPaymentType        string
	OtherActiveSubscriptionId string `json:"otherActiveSubscriptionId" description:"other active or incomplete subscription id "`
	ApplyPromoCredit          bool   `json:"applyPromoCredit" `
}

type CreateInternalReq struct {
	MerchantId             uint64                      `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId                 uint64                      `json:"planId" dc:"PlanId" v:"required"`
	UserId                 uint64                      `json:"userId" dc:"UserId" v:"required"`
	DiscountCode           string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount               *bean.ExternalDiscountParam `json:"discount" dc:"Discount"`
	Quantity               int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId              *uint64                     `json:"gatewayId" dc:"Id" `
	GatewayPaymentType     string                      `json:"gatewayPaymentType" dc:"Gateway Payment Type"`
	AddonParams            []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount     int64                       `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency        string                      `json:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl              string                      `json:"returnUrl"  dc:"RedirectUrl"  `
	CancelUrl              string                      `json:"cancelUrl" dc:"CancelUrl"`
	VatCountryCode         string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber              string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage          *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	PaymentMethodId        string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata               map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	TrialEnd               int64                       `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	StartIncomplete        bool                        `json:"StartIncomplete"        dc:"StartIncomplete, use now pay later, subscription will generate invoice and start with incomplete status if set"`
	ProductData            *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
	ApplyPromoCredit       bool                        `json:"applyPromoCredit" `
	ApplyPromoCreditAmount *int64                      `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
}

type CreateInternalRes struct {
	Plan         *entity.Plan       `json:"plan"`
	Subscription *bean.Subscription `json:"subscription" dc:"Subscription"`
	User         *bean.UserAccount  `json:"user" dc:"user"`
	Paid         bool               `json:"paid"`
	Link         string             `json:"link"`
}

func SubscriptionCreatePreview(ctx context.Context, req *CreatePreviewInternalReq) (*CreatePreviewInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	if req.IsSubmit {
		utility.Assert(req.UserId > 0, "UserId invalid")
		utility.Assert(req.GatewayId != nil, "Gateway invalid")
	}
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v not active", plan.Id))
	utility.Assert(plan.Type != consts.PlanTypeRecurringAddon, fmt.Sprintf("Plan Id:%v is addon", plan.Id))
	var user *entity.UserAccount = nil
	if req.UserId > 0 || req.IsSubmit {
		user = query.GetUserAccountById(ctx, req.UserId)
		utility.Assert(user != nil, "user not found")
	}
	var gatewayId uint64 = 0
	if req.GatewayId != nil {
		gatewayId = *req.GatewayId
	}
	var paymentMethodId = req.PaymentMethodId
	var paymentType = req.GatewayPaymentType
	if user != nil {
		gatewayId, paymentType, paymentMethodId = sub_update.VerifyPaymentGatewayMethod(ctx, user.Id, req.GatewayId, req.GatewayPaymentType, req.PaymentMethodId, "")
	}
	var gateway *entity.MerchantGateway
	if gatewayId > 0 || req.IsSubmit {
		utility.Assert(gatewayId > 0, "gateway need specified")
		gateway = query.GetGatewayById(ctx, gatewayId)
		utility.Assert(gateway != nil, "gateway not found")
		utility.Assert(gateway.MerchantId == req.MerchantId, "invalid gateway")
	}
	if !_interface.Context().Get(ctx).IsOpenApiCall && user != nil && gatewayId > 0 {
		sub_update.UpdateUserDefaultGatewayPaymentMethod(ctx, user.Id, gatewayId, paymentMethodId, req.GatewayPaymentType)
	}
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	req.Quantity = utility.MaxInt64(1, req.Quantity)
	userEmail := ""
	if user != nil {
		userEmail = user.Email
	}

	var err error
	var otherSameProductActiveSubscription *entity.Subscription
	var otherActiveSubscriptionId string
	if user != nil {
		otherSameProductActiveSubscription = query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, user.Id, merchantInfo.Id, plan.ProductId)
		if otherSameProductActiveSubscription != nil {
			otherActiveSubscriptionId = otherSameProductActiveSubscription.SubscriptionId
		}
		if req.IsSubmit {
			utility.Assert(otherSameProductActiveSubscription == nil, i18n.LocalizationFormat(ctx, "{#SubDuplicateCreation}"))
		}
	}

	var vatCountryCode = req.VatCountryCode
	var subscriptionTaxPercentage int64 = 0
	var vatCountryName = ""
	var vatCountryRate *bean.VatCountryRate
	var vatNumberValidate *bean.ValidResult
	var vatNumberValidateMessage string
	var recurringDiscountCode string
	var discountMessage string

	if len(req.VatNumber) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id).VatRatesEnabled(), i18n.LocalizationFormat(ctx, "{#VatGatewayNeedSetup}"))
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.UserId, req.VatNumber, "")
		if err != nil || !vatNumberValidate.Valid {
			if err != nil {
				g.Log().Error(ctx, "ValidateVatNumberByDefaultGateway error:%s", err.Error())
				vatNumberValidateMessage = "Server Error"
			} else {
				vatNumberValidateMessage = i18n.LocalizationFormat(ctx, "{#VatValidateError}", req.VatNumber)
			}
		} else {
			if len(req.VatCountryCode) > 0 {
				utility.Assert(vatCountryCode == vatNumberValidate.CountryCode, i18n.LocalizationFormat(ctx, "{#CountryCodeVatNumberNotMatch}", vatNumberValidate.CountryCode))
			}
			vatCountryCode = vatNumberValidate.CountryCode
		}
		if req.IsSubmit {
			utility.Assert(vatNumberValidate != nil && vatNumberValidate.Valid, i18n.LocalizationFormat(ctx, "{#VatValidateError}", req.VatNumber))
		}
	}

	var validVatNumber = ""
	if vatNumberValidate != nil && vatNumberValidate.Valid {
		validVatNumber = vatNumberValidate.VatNumber
	}
	if req.TaxPercentage != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "External TaxPercentage only available for api call")
		utility.Assert(*req.TaxPercentage >= 0 && *req.TaxPercentage < 10000, "invalid taxPercentage")
		subscriptionTaxPercentage = *req.TaxPercentage
	} else if len(vatCountryCode) > 0 && gateway != nil {
		utility.Assert(service2.IsGatewaySupportCountryCode(ctx, gateway, req.VatCountryCode), "gateway not support countryCode:"+vatCountryCode)
		taxPercentage, _ := vat_gateway.ComputeMerchantVatPercentage(ctx, req.MerchantId, vatCountryCode, gateway.Id, validVatNumber)
		subscriptionTaxPercentage = taxPercentage
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

	promoCreditDiscountCodeExclusive := config.CheckCreditConfigDiscountCodeExclusive(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, plan.Currency)
	if len(req.DiscountCode) > 0 {
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:         req.MerchantId,
			UserId:             req.UserId,
			DiscountCode:       req.DiscountCode,
			Currency:           plan.Currency,
			PLanId:             plan.Id,
			TimeNow:            gtime.Now().Timestamp(),
			IsUpgrade:          false,
			IsChangeToLongPlan: false,
			IsRenew:            false,
			IsNewUser:          IsNewSubscriptionUser(ctx, req.MerchantId, userEmail),
		})
		if canApply {
			if isRecurring {
				recurringDiscountCode = req.DiscountCode
			}
		} else {
			req.DiscountCode = ""
			discountMessage = message
		}
		{
			//conflict, disable discount code
			if promoCreditDiscountCodeExclusive && canApply && req.ApplyPromoCredit != nil && *req.ApplyPromoCredit {
				_, promoCreditPayout, _ := payment.CheckCreditUserPayout(ctx, req.MerchantId, req.UserId, consts.CreditAccountTypePromo, plan.Currency, plan.Amount, req.ApplyPromoCreditAmount)
				if promoCreditPayout != nil && promoCreditPayout.CurrencyAmount > 0 {
					discountMessage = "Promo Credit Conflict with Discount code"
					req.DiscountCode = ""
					if req.IsSubmit {
						utility.Assert(false, discountMessage)
					}
				}
			}
		}
		if req.IsSubmit {
			utility.Assert(canApply, message)
		}
	}

	if req.ApplyPromoCredit == nil {
		if promoCreditDiscountCodeExclusive && len(req.DiscountCode) > 0 {
			req.ApplyPromoCredit = unibee.Bool(false)
		} else {
			req.ApplyPromoCredit = unibee.Bool(config.CheckCreditConfigPreviewDefaultUsed(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, plan.Currency))
		}
	}

	var currentTimeStart = gtime.Now()
	var trialEnd = currentTimeStart.Timestamp() - 1
	var cancelAtPeriodEnd = 0
	utility.Assert(len(plan.IntervalUnit) > 0, "Invalid plan billing period")
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
		trialEnd = req.TrialEnd
		//var currentTimeEnd = req.TrialEnd
		var currentTimeEnd = subscription2.GetPeriodEndFromStart(ctx, currentTimeStart.Timestamp(), currentTimeStart.Timestamp(), req.PlanId)
		//Promo Credit
		var promoCreditDiscountAmount int64 = 0
		var promoCreditAccount *bean.CreditAccount
		var promoCreditPayout *bean.CreditPayout
		var creditPayoutErr error
		if *req.ApplyPromoCredit {
			promoCreditAccount, promoCreditPayout, creditPayoutErr = payment.CheckCreditUserPayout(ctx, req.MerchantId, req.UserId, consts.CreditAccountTypePromo, plan.Currency, totalAmountExcludingTax, req.ApplyPromoCreditAmount)
			if creditPayoutErr == nil && promoCreditAccount != nil && promoCreditPayout != nil {
				promoCreditDiscountAmount = promoCreditPayout.CurrencyAmount
				totalAmountExcludingTax = totalAmountExcludingTax - promoCreditDiscountAmount
			}
		}

		discountAmount := utility.MinInt64(discount.ComputeDiscountAmount(ctx, plan.MerchantId, totalAmountExcludingTax, plan.Currency, req.DiscountCode, currentTimeStart.Timestamp()), totalAmountExcludingTax)
		totalAmountExcludingTax = totalAmountExcludingTax - discountAmount

		var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(subscriptionTaxPercentage))
		var name = plan.PlanName
		var description = plan.Description
		if req.ProductData != nil && len(req.ProductData.Name) > 0 {
			name = req.ProductData.Name
			description = req.ProductData.Description
		}
		invoice := &bean.Invoice{
			InvoiceName:                    "SubscriptionCreate",
			ProductName:                    name,
			OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount + promoCreditDiscountAmount,
			TotalAmount:                    totalAmountExcludingTax + taxAmount,
			TotalAmountExcludingTax:        totalAmountExcludingTax,
			DiscountCode:                   req.DiscountCode,
			DiscountAmount:                 discountAmount,
			PromoCreditDiscountAmount:      promoCreditDiscountAmount,
			PromoCreditAccount:             promoCreditAccount,
			PromoCreditPayout:              promoCreditPayout,
			Currency:                       plan.Currency,
			TaxAmount:                      taxAmount,
			BizType:                        consts.BizTypeSubscription,
			SubscriptionAmount:             totalAmountExcludingTax + discountAmount + promoCreditDiscountAmount + taxAmount,
			SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount + promoCreditDiscountAmount,
			TrialEnd:                       trialEnd,
			Lines: []*bean.InvoiceItemSimplify{{
				Currency:               plan.Currency,
				OriginAmount:           totalAmountExcludingTax + taxAmount + discountAmount + promoCreditDiscountAmount,
				Amount:                 totalAmountExcludingTax + taxAmount,
				DiscountAmount:         discountAmount,
				Tax:                    taxAmount,
				AmountExcludingTax:     totalAmountExcludingTax + discountAmount + promoCreditDiscountAmount,
				TaxPercentage:          subscriptionTaxPercentage,
				UnitAmountExcludingTax: plan.TrialAmount,
				Name:                   name,
				Description:            description,
				Proration:              false,
				Quantity:               req.Quantity,
				PeriodEnd:              currentTimeEnd,
				PeriodStart:            currentTimeStart.Timestamp(),
				Plan:                   bean.SimplifyPlan(plan),
			}},
			PeriodStart:        currentTimeStart.Timestamp(),
			PeriodEnd:          currentTimeEnd,
			BillingCycleAnchor: currentTimeStart.Timestamp(),
			FinishTime:         currentTimeStart.Timestamp(),
			Metadata:           req.Metadata,
			VatNumber:          validVatNumber,
			CountryCode:        vatCountryCode,
			TaxPercentage:      subscriptionTaxPercentage,
		}
		return &CreatePreviewInternalRes{
			Plan:                      plan,
			User:                      bean.SimplifyUserAccount(user),
			TrialEnd:                  trialEnd,
			Quantity:                  req.Quantity,
			Gateway:                   gateway,
			Merchant:                  merchantInfo,
			AddonParams:               req.AddonParams,
			Addons:                    addons,
			OriginAmount:              invoice.OriginAmount,
			TotalAmount:               invoice.TotalAmount,
			DiscountAmount:            invoice.DiscountAmount,
			Invoice:                   invoice,
			RecurringDiscountCode:     recurringDiscountCode,
			Discount:                  bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, invoice.DiscountCode)),
			Currency:                  currency,
			VatCountryCode:            vatCountryCode,
			VatCountryName:            vatCountryName,
			VatNumber:                 req.VatNumber,
			VatNumberValidate:         vatNumberValidate,
			VatVerifyData:             utility.MarshalToJsonString(vatNumberValidate),
			UserId:                    req.UserId,
			Email:                     userEmail,
			VatCountryRate:            vatCountryRate,
			Gateways:                  service2.GetMerchantAvailableGatewaysByCountryCode(ctx, req.MerchantId, req.VatCountryCode),
			TaxPercentage:             subscriptionTaxPercentage,
			VatNumberValidateMessage:  vatNumberValidateMessage,
			DiscountMessage:           discountMessage,
			OtherActiveSubscriptionId: otherActiveSubscriptionId,
			ApplyPromoCredit:          *req.ApplyPromoCredit,
		}, nil
	} else {
		var currentTimeEnd = subscription2.GetPeriodEndFromStart(ctx, currentTimeStart.Timestamp(), currentTimeStart.Timestamp(), req.PlanId)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			UserId:                 req.UserId,
			InvoiceName:            "SubscriptionCreate",
			DiscountCode:           req.DiscountCode,
			TimeNow:                gtime.Now().Timestamp(),
			Currency:               currency,
			PlanId:                 req.PlanId,
			Quantity:               req.Quantity,
			AddonJsonData:          utility.MarshalToJsonString(req.AddonParams),
			CountryCode:            vatCountryCode,
			VatNumber:              validVatNumber,
			TaxPercentage:          subscriptionTaxPercentage,
			PeriodStart:            currentTimeStart.Timestamp(),
			PeriodEnd:              currentTimeEnd,
			FinishTime:             currentTimeStart.Timestamp(),
			ProductData:            req.ProductData,
			BillingCycleAnchor:     currentTimeStart.Timestamp(),
			Metadata:               req.Metadata,
			ApplyPromoCredit:       *req.ApplyPromoCredit,
			ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
		})

		return &CreatePreviewInternalRes{
			Plan:                      plan,
			User:                      bean.SimplifyUserAccount(user),
			TrialEnd:                  trialEnd,
			Quantity:                  req.Quantity,
			Gateway:                   gateway,
			Merchant:                  merchantInfo,
			AddonParams:               req.AddonParams,
			Addons:                    addons,
			OriginAmount:              invoice.OriginAmount,
			TotalAmount:               invoice.TotalAmount,
			DiscountAmount:            invoice.DiscountAmount,
			Invoice:                   invoice,
			RecurringDiscountCode:     recurringDiscountCode,
			Discount:                  bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, invoice.DiscountCode)),
			Currency:                  currency,
			VatCountryCode:            vatCountryCode,
			VatCountryName:            vatCountryName,
			VatNumber:                 req.VatNumber,
			VatNumberValidate:         vatNumberValidate,
			VatVerifyData:             utility.MarshalToJsonString(vatNumberValidate),
			UserId:                    req.UserId,
			Email:                     userEmail,
			VatCountryRate:            vatCountryRate,
			Gateways:                  service2.GetMerchantAvailableGatewaysByCountryCode(ctx, req.MerchantId, req.VatCountryCode),
			TaxPercentage:             subscriptionTaxPercentage,
			VatNumberValidateMessage:  vatNumberValidateMessage,
			DiscountMessage:           discountMessage,
			CancelAtPeriodEnd:         cancelAtPeriodEnd,
			GatewayPaymentMethodId:    paymentMethodId,
			GatewayPaymentType:        paymentType,
			OtherActiveSubscriptionId: otherActiveSubscriptionId,
			ApplyPromoCredit:          *req.ApplyPromoCredit,
		}, nil
	}
}

func SubscriptionCreate(ctx context.Context, req *CreateInternalReq) (*CreateInternalRes, error) {
	if req.Discount != nil {
		// create external discount
		utility.Assert(req.PlanId > 0, "PlanId invalid")
		utility.Assert(req.UserId > 0, "UserId invalid")
		plan := query.GetPlanById(ctx, req.PlanId)
		utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
		utility.Assert(plan != nil, "invalid planId")
		one := discount.CreateExternalDiscount(ctx, req.MerchantId, req.UserId, strconv.FormatUint(req.PlanId, 10), req.Discount, plan.Currency, gtime.Now().Timestamp())
		req.DiscountCode = one.Code
	} else if len(req.DiscountCode) > 0 {
		one := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
		utility.Assert(one != nil, i18n.LocalizationFormat(ctx, "{#DiscountCodeInvalid}"))
		utility.Assert(one.Type == 0, "invalid code, code is from external")
	}

	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	existOne := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, req.UserId, req.MerchantId, plan.ProductId)
	if existOne != nil {
		err := SubscriptionCancel(ctx, existOne.SubscriptionId, false, false, "CancelledByAnotherCreation")
		utility.AssertError(err, "Subscription cancel error")
	}

	prepare, err := SubscriptionCreatePreview(ctx, &CreatePreviewInternalReq{
		MerchantId:             req.MerchantId,
		PlanId:                 req.PlanId,
		UserId:                 req.UserId,
		DiscountCode:           req.DiscountCode,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		AddonParams:            req.AddonParams,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		TaxPercentage:          req.TaxPercentage,
		IsSubmit:               true,
		TrialEnd:               req.TrialEnd,
		ProductData:            req.ProductData,
		PaymentMethodId:        req.PaymentMethodId,
		Metadata:               req.Metadata,
		ApplyPromoCredit:       unibee.Bool(req.ApplyPromoCredit),
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	})
	if err != nil {
		return nil, err
	}
	// todo mark countryCode is required or node
	// utility.Assert(len(prepare.VatCountryCode) > 0, "CountryCode Needed")
	if req.ConfirmTotalAmount > 0 {
		utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, i18n.LocalizationFormat(ctx, "{#AmountNotMatch}"))
	}
	if len(req.ConfirmCurrency) > 0 {
		utility.Assert(strings.Compare(strings.ToUpper(req.ConfirmCurrency), prepare.Currency) == 0, "currency not match , data may expired, fetch preview again")
	}

	if prepare.Invoice.TotalAmount == 0 && strings.Contains(prepare.Plan.TrialDemand, "paymentMethod") && req.PaymentMethodId == "" {
		utility.Assert(prepare.Gateway.GatewayType == consts.GatewayTypeCard || prepare.Gateway.GatewayType == consts.GatewayTypePaypal, i18n.LocalizationFormat(ctx, "{#PlanTrialGatewayError}"))
	}

	var subType = consts.SubTypeDefault
	if consts.SubscriptionCycleUnderUniBeeControl {
		subType = consts.SubTypeUniBeeControl
	}

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(prepare.Invoice.PeriodEnd, prepare.TrialEnd), prepare.Plan.Id)

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
		BillingCycleAnchor:          prepare.Invoice.BillingCycleAnchor,
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

	{
		if !vat_gateway.GetDefaultVatGateway(ctx, req.MerchantId).VatRatesEnabled() {
			sub_update.UpdateUserTaxPercentageOnly(ctx, prepare.UserId, one.TaxPercentage)
		}
	}

	var createRes *gateway_bean.GatewayCreateSubscriptionResp
	invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
		PlanId:             req.PlanId,
		Simplify:           prepare.Invoice,
		Sub:                one,
		GatewayId:          one.GatewayId,
		GatewayPaymentType: prepare.GatewayPaymentType,
		PaymentMethodId:    prepare.GatewayPaymentMethodId,
		IsSubLatestInvoice: true,
		TimeNow:            gtime.Now().Timestamp(),
	})
	utility.AssertError(err, "System Error")
	timeline.SubscriptionNewPendingTimeline(ctx, invoice)
	if prepare.Invoice.TotalAmount == 0 {
		//totalAmount is 0, no payment need
		utility.AssertError(err, "System Error")
		if strings.Contains(prepare.Plan.TrialDemand, "paymentMethod") && req.PaymentMethodId == "" {
			url, _ := method.NewPaymentMethod(ctx, &method.NewPaymentMethodInternalReq{
				MerchantId:     _interface.GetMerchantId(ctx),
				UserId:         req.UserId,
				Currency:       prepare.Currency,
				GatewayId:      prepare.Gateway.Id,
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
			sub_update.UpdateUserCountryCode(ctx, one.UserId, one.CountryCode)
			createRes = &gateway_bean.GatewayCreateSubscriptionResp{
				GatewaySubscriptionId: one.SubscriptionId,
				Link:                  GetSubscriptionZeroPaymentLink(req.ReturnUrl, one.SubscriptionId),
				Paid:                  true,
			}
		}
	} else {
		createPaymentResult, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, &service.CreateSubInvoicePaymentDefaultAutomaticReq{
			Invoice:       invoice,
			ManualPayment: len(req.PaymentMethodId) == 0,
			ReturnUrl:     req.ReturnUrl,
			CancelUrl:     req.CancelUrl,
			Source:        "SubscriptionCreate",
			TimeNow:       0,
		})
		if err != nil {
			// todo mark use method
			_, updateErr := dao.Subscription.Ctx(ctx).Data(g.Map{
				dao.Subscription.Columns().Status:       consts.SubStatusCancelled,
				dao.Subscription.Columns().CancelReason: "Create First Payment Error",
				dao.Subscription.Columns().GmtModify:    gtime.Now(),
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
		Topic:      redismq2.TopicSubscriptionCreate.Topic,
		Tag:        redismq2.TopicSubscriptionCreate.Tag,
		Body:       one.SubscriptionId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
		Content:        fmt.Sprintf("Creation(%s)", consts.SubStatusToEnum(one.Status).Description()),
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      invoice.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
	if req.GatewayId != nil {
		sub_update.UpdateUserDefaultGatewayForCheckout(ctx, req.UserId, *req.GatewayId, req.GatewayPaymentType)
	}
	return &CreateInternalRes{
		Plan:         prepare.Plan,
		Subscription: bean.SimplifySubscription(ctx, one),
		User:         prepare.User,
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
}
