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
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/email"
	"unibee/internal/logic/gateway/gateway_bean"
	service2 "unibee/internal/logic/gateway/service"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/method"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/config"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/subscription/pending_update_cancel"
	"unibee/internal/logic/subscription/timeline"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/vat_gateway"
	entity "unibee/internal/model/entity/default"
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

func VatNumberValidate(ctx context.Context, req *vat.NumberValidateReq) (*vat.NumberValidateRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.VatNumber) > 0, "vatNumber invalid")
	vatNumberValidate, err := vat_gateway.ValidateVatNumberByDefaultGateway(ctx, _interface.GetMerchantId(ctx), 0, req.VatNumber, "")
	if err != nil {
		return nil, err
	}
	//if vatNumberValidate.Valid {
	//	vatCountryRate, err := vat_gateway.QueryVatCountryRateByMerchant(ctx, _interface.GetMerchantId(ctx), vatNumberValidate.CountryCode)
	//	utility.Assert(err == nil, fmt.Sprintf("verify error:%s", err))
	//	utility.Assert(vatCountryRate != nil, fmt.Sprintf("vatNumber not found for countryCode:%v", vatNumberValidate.CountryCode))
	//}
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
	CancelUrl     string                      `json:"cancelUrl" dc:"CancelUrl"`
	ProductData   *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
	Metadata      map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
}

func GetSubscriptionZeroPaymentLink(returnUrl string, subId string) string {
	if returnUrl == "" {
		return returnUrl
	}
	if returnUrl != "" && strings.Contains(returnUrl, "?") {
		return fmt.Sprintf("%s&subId=%s&success=true", returnUrl, subId)
	} else {
		return fmt.Sprintf("%s?subId=%s&success=true", returnUrl, subId)
	}
}

func SubscriptionRenew(ctx context.Context, req *RenewInternalReq) (*CreateInternalRes, error) {
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.MerchantId == req.MerchantId, "merchantId not match")
	// todo mark renew for all status
	//utility.Assert(sub.Status == consts.SubStatusExpired || sub.Status == consts.SubStatusCancelled, "subscription not cancel or expire status")
	var subscriptionTaxPercentage = sub.TaxPercentage
	percentage, countryCode, vatNumber, err := sub_update.GetUserTaxPercentage(ctx, sub.UserId)
	if err == nil {
		subscriptionTaxPercentage = percentage
	}
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
	//var gatewayId = sub.GatewayId
	//if req.GatewayId != nil {
	//	gatewayId = *req.GatewayId
	//}
	gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, req.GatewayId, "", sub.SubscriptionId)
	utility.Assert(gatewayId > 0, "gateway need specified")
	var timeNow = gtime.Now().Timestamp()
	if sub.TestClock > timeNow && !config2.GetConfigInstance().IsProd() {
		timeNow = sub.TestClock
	}

	if req.Discount != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		// create external discount
		utility.Assert(sub.PlanId > 0, "planId invalid")
		plan := query.GetPlanById(ctx, sub.PlanId)
		utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
		utility.Assert(plan != nil, "invalid planId")
		one := discount.CreateExternalDiscount(ctx, req.MerchantId, sub.UserId, strconv.FormatUint(sub.PlanId, 10), req.Discount, plan.Currency, utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock))
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
			TimeNow:        utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
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
			TimeNow:        utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
		})
		if canApply && isRecurring {
			req.DiscountCode = sub.DiscountCode
		}
	}

	currentInvoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:        "SubscriptionRenew",
		Currency:           sub.Currency,
		DiscountCode:       req.DiscountCode,
		TimeNow:            timeNow,
		PlanId:             sub.PlanId,
		Quantity:           sub.Quantity,
		AddonJsonData:      utility.MarshalToJsonString(addonParams),
		CountryCode:        countryCode,
		VatNumber:          vatNumber,
		TaxPercentage:      subscriptionTaxPercentage,
		PeriodStart:        timeNow,
		PeriodEnd:          subscription2.GetPeriodEndFromStart(ctx, timeNow, timeNow, sub.PlanId),
		FinishTime:         timeNow,
		ProductData:        req.ProductData,
		BillingCycleAnchor: timeNow,
		Metadata:           req.Metadata,
	})

	// createAndPayNewProrationInvoice
	merchantInfo := query.GetMerchantById(ctx, sub.MerchantId)
	utility.Assert(merchantInfo != nil, "merchantInfo not found")
	// utility.Assert(user != nil, "user not found")
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	invoice, err := service3.CreateProcessingInvoiceForSub(ctx, currentInvoice, sub, gateway.Id, paymentMethodId, true, timeNow)
	utility.AssertError(err, "System Error")
	var createRes *gateway_bean.GatewayNewPaymentResp
	if invoice.TotalAmount > 0 {
		createRes, err = service.CreateSubInvoicePaymentDefaultAutomatic(ctx, invoice, req.ManualPayment, req.ReturnUrl, req.CancelUrl, "SubscriptionRenew", 0)
		if err != nil {
			g.Log().Print(ctx, "SubscriptionRenew CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return nil, err
		}
	} else {
		invoice, err = handler2.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
		utility.AssertError(err, "System Error")
		createRes = &gateway_bean.GatewayNewPaymentResp{
			Payment:                nil,
			Status:                 consts.PaymentSuccess,
			GatewayPaymentId:       "",
			GatewayPaymentIntentId: "",
			GatewayPaymentMethod:   "",
			Link:                   GetSubscriptionZeroPaymentLink(req.ReturnUrl, sub.SubscriptionId),
			Action:                 nil,
			Invoice:                nil,
			PaymentCode:            "",
		}
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", sub.SubscriptionId),
		Content:        "Renew",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      invoice.InvoiceId,
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	// need cancel payment、 invoice and send invoice email
	pending_update_cancel.CancelOtherUnfinishedPendingUpdatesBackground(sub.SubscriptionId, sub.SubscriptionId, "CancelByRenewSubscription-"+sub.SubscriptionId)

	if createRes.Status == consts.PaymentSuccess {
		err = handler.HandleSubscriptionNextBillingCyclePaymentSuccess(ctx, sub, invoice)
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
		Paid:         createRes.Status == consts.PaymentSuccess,
		Link:         createRes.Link,
	}, nil
}

type CreatePreviewInternalReq struct {
	MerchantId      uint64                 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId          uint64                 `json:"planId" dc:"PlanId" v:"required"`
	UserId          uint64                 `json:"userId" dc:"UserId" v:"required"`
	Quantity        int64                  `json:"quantity" dc:"Quantity" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
	GatewayId       *uint64                `json:"gatewayId" dc:"Id"`
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	VatCountryCode  string                 `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber       string                 `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage   *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	TrialEnd        int64                  `json:"trialEnd"  description:"trial_end, utc time"` // trial_end, utc time
	IsSubmit        bool
	ProductData     *bean.PlanProductParam `json:"productData"  dc:"ProductData"  `
	PaymentMethodId string
	Metadata        map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type CreatePreviewInternalRes struct {
	Plan                     *entity.Plan               `json:"plan"`
	User                     *bean.UserAccount          `json:"user"`
	Quantity                 int64                      `json:"quantity"`
	Gateway                  *entity.MerchantGateway    `json:"gateway"`
	Merchant                 *entity.Merchant           `json:"merchantInfo"`
	AddonParams              []*bean.PlanAddonParam     `json:"addonParams"`
	Addons                   []*bean.PlanAddonDetail    `json:"addons"`
	OriginAmount             int64                      `json:"originAmount"                `
	TotalAmount              int64                      `json:"totalAmount" `
	DiscountAmount           int64                      `json:"discountAmount"`
	Currency                 string                     `json:"currency" `
	VatCountryCode           string                     `json:"vatCountryCode" `
	VatCountryName           string                     `json:"vatCountryName" `
	VatNumber                string                     `json:"vatNumber" `
	VatNumberValidate        *bean.ValidResult          `json:"vatNumberValidate" `
	TaxPercentage            int64                      `json:"taxPercentage" `
	TrialEnd                 int64                      `json:"trialEnd" `
	VatVerifyData            string                     `json:"vatVerifyData" `
	Invoice                  *bean.Invoice              `json:"invoice"`
	UserId                   uint64                     `json:"userId" `
	Email                    string                     `json:"email" `
	VatCountryRate           *bean.VatCountryRate       `json:"vatCountryRate" `
	Gateways                 []*bean.Gateway            `json:"gateways" `
	RecurringDiscountCode    string                     `json:"recurringDiscountCode" `
	Discount                 *bean.MerchantDiscountCode `json:"discount" `
	VatNumberValidateMessage string                     `json:"vatNumberValidateMessage" `
	DiscountMessage          string                     `json:"discountMessage" `
	CancelAtPeriodEnd        int                        `json:"cancelAtPeriodEnd"           description:"whether cancel at period end，0-false | 1-true"` // whether cancel at period end，0-false | 1-true
	GatewayPaymentMethodId   string
}

type CreateInternalReq struct {
	MerchantId         uint64                      `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	PlanId             uint64                      `json:"planId" dc:"PlanId" v:"required"`
	UserId             uint64                      `json:"userId" dc:"UserId" v:"required"`
	DiscountCode       string                      `json:"discountCode"        dc:"DiscountCode"`
	Discount           *bean.ExternalDiscountParam `json:"discount" dc:"Discount"`
	Quantity           int64                       `json:"quantity" dc:"Quantity，Default 1" `
	GatewayId          *uint64                     `json:"gatewayId" dc:"Id" `
	AddonParams        []*bean.PlanAddonParam      `json:"addonParams" dc:"addonParams" `
	ConfirmTotalAmount int64                       `json:"confirmTotalAmount"  dc:"TotalAmount To Be Confirmed，Get From Preview"  v:"required"            `
	ConfirmCurrency    string                      `json:"confirmCurrency"  dc:"Currency To Be Confirmed，Get From Preview" v:"required"  `
	ReturnUrl          string                      `json:"returnUrl"  dc:"RedirectUrl"  `
	CancelUrl          string                      `json:"cancelUrl" dc:"CancelUrl"`
	VatCountryCode     string                      `json:"vatCountryCode" dc:"VatCountryCode, CountryName"`
	VatNumber          string                      `json:"vatNumber" dc:"VatNumber" `
	TaxPercentage      *int64                      `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%"`
	PaymentMethodId    string                      `json:"paymentMethodId" dc:"PaymentMethodId" `
	Metadata           map[string]interface{}      `json:"metadata" dc:"Metadata，Map"`
	TrialEnd           int64                       `json:"trialEnd"                    description:"trial_end, utc time"` // trial_end, utc time
	StartIncomplete    bool                        `json:"StartIncomplete"        dc:"StartIncomplete, use now pay later, subscription will generate invoice and start with incomplete status if set"`
	ProductData        *bean.PlanProductParam      `json:"productData"  dc:"ProductData"  `
}

type CreateInternalRes struct {
	Subscription *bean.Subscription `json:"subscription" dc:"Subscription"`
	User         *bean.UserAccount  `json:"user" dc:"user"`
	Paid         bool               `json:"paid"`
	Link         string             `json:"link"`
}

func SubscriptionCreatePreview(ctx context.Context, req *CreatePreviewInternalReq) (*CreatePreviewInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.PlanId > 0, "PlanId invalid")
	utility.Assert(req.UserId > 0, "UserId invalid")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "invalid planId")
	utility.Assert(plan.MerchantId == req.MerchantId, "merchant not match")
	utility.Assert(plan.Status == consts.PlanStatusActive, fmt.Sprintf("Plan Id:%v Not Publish status", plan.Id))
	utility.Assert(plan.Type == consts.PlanTypeMain, fmt.Sprintf("Plan Id:%v Not Main Type", plan.Id))
	user := query.GetUserAccountById(ctx, req.UserId)
	utility.Assert(user != nil, "user not found")
	gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, req.UserId, req.GatewayId, req.PaymentMethodId, "")
	utility.Assert(gatewayId > 0, "gateway need specified")
	if !_interface.Context().Get(ctx).IsOpenApiCall {
		sub_update.UpdateUserDefaultGatewayPaymentMethod(ctx, user.Id, gatewayId, paymentMethodId)
	}
	gateway := MerchantGatewayCheck(ctx, plan.MerchantId, gatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	utility.Assert(service2.IsGatewaySupportCountryCode(ctx, gateway, req.VatCountryCode), "gateway not support")
	merchantInfo := query.GetMerchantById(ctx, plan.MerchantId)
	utility.Assert(merchantInfo != nil, "merchant not found")

	req.Quantity = utility.MaxInt64(1, req.Quantity)

	var err error
	utility.Assert(query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, merchantInfo.Id, plan.ProductId) == nil, "Another active or incomplete subscription exist")

	//setup vat from user
	//if len(req.VatCountryCode) == 0 && len(user.CountryCode) > 0 {
	//	req.VatCountryCode = user.CountryCode
	//}
	//if len(req.VatNumber) == 0 {
	//	req.VatNumber = user.VATNumber
	//}

	var vatCountryCode = req.VatCountryCode
	var subscriptionTaxPercentage int64 = 0
	var vatCountryName = ""
	var vatCountryRate *bean.VatCountryRate
	var vatNumberValidate *bean.ValidResult
	var vatNumberValidateMessage string
	var discountMessage string

	//if len(req.VatCountryCode) == 0 {
	//	req.VatNumber = user.CountryCode
	//}

	if len(req.VatNumber) > 0 {
		utility.Assert(vat_gateway.GetDefaultVatGateway(ctx, merchantInfo.Id) != nil, "Vat gateway need setup")
		vatNumberValidate, err = vat_gateway.ValidateVatNumberByDefaultGateway(ctx, merchantInfo.Id, req.UserId, req.VatNumber, "")
		if err != nil || !vatNumberValidate.Valid {
			if err != nil {
				g.Log().Error(ctx, "ValidateVatNumberByDefaultGateway error:%s", err.Error())
				vatNumberValidateMessage = "Server Error"
			} else {
				vatNumberValidateMessage = "Validate Failure"
			}
		} else {
			if len(req.VatCountryCode) > 0 && !_interface.Context().Get(ctx).IsOpenApiCall {
				utility.Assert(vatCountryCode == vatNumberValidate.CountryCode, "CountryCode error, "+"Your country from vat number is "+vatNumberValidate.CountryCode)
			}
			vatCountryCode = vatNumberValidate.CountryCode
		}
		if req.IsSubmit {
			utility.Assert(vatNumberValidate.Valid, fmt.Sprintf("VatNumber validate failure, number:"+req.VatNumber))
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
	} else if len(vatCountryCode) > 0 {
		taxPercentage, _ := vat_gateway.ComputeMerchantVatPercentage(ctx, user.MerchantId, vatCountryCode, gatewayId, validVatNumber)
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

	var recurringDiscountCode string
	if len(req.DiscountCode) > 0 {
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:   req.MerchantId,
			UserId:       req.UserId,
			DiscountCode: req.DiscountCode,
			Currency:     plan.Currency,
			PLanId:       plan.Id,
			TimeNow:      gtime.Now().Timestamp(),
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
		var name = plan.PlanName
		var description = plan.Description
		if req.ProductData != nil && len(req.ProductData.Name) > 0 {
			name = req.ProductData.Name
			description = req.ProductData.Description
		}
		invoice := &bean.Invoice{
			InvoiceName:                    "SubscriptionCreate",
			ProductName:                    name,
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
			Plan:                     plan,
			User:                     bean.SimplifyUserAccount(user),
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
		var currentTimeEnd = subscription2.GetPeriodEndFromStart(ctx, currentTimeStart.Timestamp(), currentTimeStart.Timestamp(), req.PlanId)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			InvoiceName:        "SubscriptionCreate",
			DiscountCode:       req.DiscountCode,
			TimeNow:            gtime.Now().Timestamp(),
			Currency:           currency,
			PlanId:             req.PlanId,
			Quantity:           req.Quantity,
			AddonJsonData:      utility.MarshalToJsonString(req.AddonParams),
			CountryCode:        vatCountryCode,
			VatNumber:          validVatNumber,
			TaxPercentage:      subscriptionTaxPercentage,
			PeriodStart:        currentTimeStart.Timestamp(),
			PeriodEnd:          currentTimeEnd,
			FinishTime:         currentTimeStart.Timestamp(),
			ProductData:        req.ProductData,
			BillingCycleAnchor: currentTimeStart.Timestamp(),
			Metadata:           req.Metadata,
		})

		return &CreatePreviewInternalRes{
			Plan:                     plan,
			User:                     bean.SimplifyUserAccount(user),
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
			GatewayPaymentMethodId:   paymentMethodId,
		}, nil
	}
}

func SubscriptionCreate(ctx context.Context, req *CreateInternalReq) (*CreateInternalRes, error) {
	if req.Discount != nil {
		//utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call") // todo mark enable for test automatic
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
		MerchantId:      req.MerchantId,
		PlanId:          req.PlanId,
		UserId:          req.UserId,
		DiscountCode:    req.DiscountCode,
		Quantity:        req.Quantity,
		GatewayId:       req.GatewayId,
		AddonParams:     req.AddonParams,
		VatCountryCode:  req.VatCountryCode,
		VatNumber:       req.VatNumber,
		TaxPercentage:   req.TaxPercentage,
		IsSubmit:        true,
		TrialEnd:        req.TrialEnd,
		ProductData:     req.ProductData,
		PaymentMethodId: req.PaymentMethodId,
		Metadata:        req.Metadata,
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
		utility.Assert(prepare.Gateway.GatewayType == consts.GatewayTypeCard, "card payment gateway need") // todo mark
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

	var createRes *gateway_bean.GatewayCreateSubscriptionResp
	invoice, err := service3.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, one, one.GatewayId, prepare.GatewayPaymentMethodId, true, gtime.Now().Timestamp())
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
			createRes = &gateway_bean.GatewayCreateSubscriptionResp{
				GatewaySubscriptionId: one.SubscriptionId,
				Link:                  GetSubscriptionZeroPaymentLink(req.ReturnUrl, one.SubscriptionId),
				Paid:                  true,
			}
		}
		// todo mark subscription become active with payment mq message
		//} else if len(req.PaymentMethodId) > 0 {
		//	// createAndPayNewProrationInvoice
		//	var createPaymentResult, err = service.CreateSubInvoicePaymentDefaultAutomatic(ctx, invoice, false, req.ReturnUrl, "SubscriptionCreate")
		//	if err != nil {
		//		g.Log().Print(ctx, "SubscriptionCreate CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
		//		return nil, err
		//	}
		//	createRes = &gateway_bean.GatewayCreateSubscriptionResp{
		//		GatewaySubscriptionId: createPaymentResult.Payment.PaymentId,
		//		Data:                  utility.MarshalToJsonString(createPaymentResult),
		//		Link:                  createPaymentResult.Link,
		//		Paid:                  createPaymentResult.Status == consts.PaymentSuccess,
		//	}
	} else {
		createPaymentResult, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, invoice, len(req.PaymentMethodId) == 0, req.ReturnUrl, req.CancelUrl, "SubscriptionCreate", 0)
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
		Topic: redismq2.TopicSubscriptionCreate.Topic,
		Tag:   redismq2.TopicSubscriptionCreate.Tag,
		Body:  one.SubscriptionId,
	})
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
		Content:        "New",
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
		if prepare.Invoice.TotalAmount == 0 {
			// zero invoice will not create payment
			_, _ = redismq.Send(&redismq.Message{
				Topic: redismq2.TopicSubscriptionPaymentSuccess.Topic,
				Tag:   redismq2.TopicSubscriptionPaymentSuccess.Tag,
				Body:  one.SubscriptionId,
			})
		}
	} else if req.StartIncomplete {
		err = SubscriptionActiveTemporarily(ctx, one.SubscriptionId, one.CurrentPeriodEnd)
		utility.AssertError(err, "Start Active Temporarily")
	}
	return &CreateInternalRes{
		Subscription: bean.SimplifySubscription(one),
		User:         prepare.User,
		Paid:         createRes.Paid,
		Link:         one.Link,
	}, nil
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
	GatewayId       *uint64                `json:"gatewayId" dc:"Id" `
	EffectImmediate int                    `json:"effectImmediate" dc:"Effect Immediate，1-Immediate，2-Next Period" `
	AddonParams     []*bean.PlanAddonParam `json:"addonParams" dc:"addonParams" `
	DiscountCode    string                 `json:"discountCode"        dc:"DiscountCode"`
	TaxPercentage   *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, override subscription taxPercentage if provide"`
	ProductData     *bean.PlanProductParam `json:"productData"  dc:"ProductData"  `
	PaymentMethodId string
	Metadata        map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
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
	utility.Assert(sub.CancelAtPeriodEnd == 0, "subscription will cancel at period end, should resume subscription first")
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

	var currentInvoice *bean.Invoice
	var nextPeriodInvoice *bean.Invoice
	var RecurringDiscountCode string
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
			TimeNow:        utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
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
				InvoiceName:        "SubscriptionUpgrade",
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
				Metadata:                       req.Metadata,
				CountryCode:                    countryCode,
				VatNumber:                      vatNumber,
				TaxPercentage:                  subscriptionTaxPercentage,
			}
		} else if prorationDate > sub.CurrentPeriodEnd {
			// after periodEnd, is not a currentInvoice, just use it
			currentInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
				InvoiceName:        "SubscriptionUpgrade",
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
			if !hasIntervalChange {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToFixedEndInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:       "SubscriptionUpgrade",
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
				})
			} else {
				currentInvoice = invoice_compute.ComputeSubscriptionProrationToDifferentIntervalInvoiceDetailSimplify(ctx, &invoice_compute.CalculateProrationInvoiceReq{
					InvoiceName:        "SubscriptionUpgrade",
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
				})
			}
		}
		prorationDate = currentInvoice.ProrationDate
	} else {
		prorationDate = utility.MaxInt64(sub.CurrentPeriodEnd, sub.TrialEnd)
		currentInvoice = &bean.Invoice{
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
			Metadata:                       req.Metadata,
			CountryCode:                    countryCode,
			VatNumber:                      vatNumber,
			TaxPercentage:                  subscriptionTaxPercentage,
		}
	}

	nextPeriodInvoice = invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
		InvoiceName:        "SubscriptionCycle",
		Currency:           sub.Currency,
		DiscountCode:       req.DiscountCode,
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
		RecurringDiscountCode: RecurringDiscountCode,
		Discount:              bean.SimplifyMerchantDiscountCode(query.GetDiscountByCode(ctx, plan.MerchantId, currentInvoice.DiscountCode)),
		PaymentMethodId:       paymentMethodId,
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

func SubscriptionUpdate(ctx context.Context, req *UpdateInternalReq, merchantMemberId int64) (*subscription.UpdateRes, error) {
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
		utility.Assert(req.ConfirmTotalAmount == prepare.TotalAmount, "totalAmount not match , data may expired, fetch preview again")
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
		invoice, err := service3.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, prepare.Subscription, prepare.Gateway.Id, prepare.PaymentMethodId, false, prepare.ProrationDate)
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
		invoice, err := service3.CreateProcessingInvoiceForSub(ctx, prepare.Invoice, prepare.Subscription, prepare.Gateway.Id, prepare.PaymentMethodId, false, prepare.ProrationDate)
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
		dao.Subscription.Columns().DiscountCode:    prepare.RecurringDiscountCode,
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
			MerchantMember:  detail.ConvertMemberToDetail(ctx, query.GetMerchantMemberById(ctx, uint64(one.MerchantMemberId))),
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
		dao.Subscription.Columns().Status:         nextStatus,
		dao.Subscription.Columns().CancelReason:   reason,
		dao.Subscription.Columns().TrialEnd:       sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()
	if err != nil {
		return err
	}
	service3.TryCancelSubscriptionLatestInvoice(ctx, sub)

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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "Cancel",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionCancelAtPeriodEnd(ctx context.Context, subscriptionId string, proration bool, merchantMemberId int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusActive, "subscription not in active status")
	if sub.CancelAtPeriodEnd == 1 {
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "EnableCancelAtPeriodEnd",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "DisableCancelAtPeriodEnd",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

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

	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd), sub.PlanId)
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
		dao.Subscription.Columns().BillingCycleAnchor: newTrialEnd,
		dao.Subscription.Columns().DunningTime:        dunningTime,
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "AddNewTrialEnd",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	return nil
}

func SubscriptionActiveTemporarily(ctx context.Context, subscriptionId string, expireTime int64) error {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing, "subscription not in pending or processing status")
	utility.Assert(sub.CurrentPeriodStart < expireTime, "expireTime should greater then subscription's period start time")
	utility.Assert(sub.CurrentPeriodEnd >= expireTime, "expireTime should lower then subscription's period end time")
	utility.Assert(len(sub.LatestInvoiceId) > 0, "sub latest invoice not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
	utility.Assert(invoice != nil, "sub latest invoice not found")
	utility.Assert(invoice.Status == consts.InvoiceStatusProcessing, "sub latest invoice not in processing")
	_, err := dao.Subscription.Ctx(ctx).Data(g.Map{
		dao.Subscription.Columns().CurrentPeriodPaid: expireTime,
		dao.Subscription.Columns().GmtModify:         gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	utility.AssertError(err, "Subscription Active Temporarily")

	_, err = dao.Invoice.Ctx(ctx).Data(g.Map{
		dao.Invoice.Columns().DayUtilDue: ((expireTime - invoice.FinishTime) / 86400) + 3,
		dao.Invoice.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.Invoice.Columns().InvoiceId, invoice.InvoiceId).OmitNil().Update()
	utility.AssertError(err, "Subscription Active Temporarily")

	err = handler.MakeSubscriptionIncomplete(ctx, subscriptionId)
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
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%v)", sub.SubscriptionId),
		Content:        "ActiveTemporarily",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)

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

func EndTrialManual(ctx context.Context, subscriptionId string) (err error) {
	utility.Assert(len(subscriptionId) > 0, "subscriptionId is nil")
	sub := query.GetSubscriptionBySubscriptionId(ctx, subscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.TrialEnd > gtime.Now().Timestamp(), "subscription not in trial period")
	newTrialEnd := sub.CurrentPeriodStart - 1
	var dunningTime = subscription2.GetDunningTimeFromEnd(ctx, utility.MaxInt64(newTrialEnd, sub.CurrentPeriodEnd), sub.PlanId)
	newStatus := sub.Status
	if gtime.Now().Timestamp() > sub.CurrentPeriodEnd {
		// todo mark has unfinished pending update
		newStatus = consts.SubStatusIncomplete
		// Payment Pending Enter Incomplete
		plan := query.GetPlanById(ctx, sub.PlanId)

		var nextPeriodStart = gtime.Now().Timestamp()
		var nextPeriodEnd = subscription2.GetPeriodEndFromStart(ctx, nextPeriodStart, nextPeriodStart, plan.Id)
		invoice := invoice_compute.ComputeSubscriptionBillingCycleInvoiceDetailSimplify(ctx, &invoice_compute.CalculateInvoiceReq{
			Currency:           sub.Currency,
			PlanId:             sub.PlanId,
			Quantity:           sub.Quantity,
			AddonJsonData:      sub.AddonData,
			TaxPercentage:      sub.TaxPercentage,
			PeriodStart:        nextPeriodStart,
			PeriodEnd:          nextPeriodEnd,
			InvoiceName:        "SubscriptionCycle",
			FinishTime:         nextPeriodStart,
			BillingCycleAnchor: nextPeriodStart,
			VatNumber:          sub.VatNumber,
		})
		gatewayId, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, nil, "", sub.SubscriptionId)
		utility.Assert(gatewayId > 0, "gateway need specified")
		one, err := service3.CreateProcessingInvoiceForSub(ctx, invoice, sub, gatewayId, paymentMethodId, true, gtime.Now().Timestamp())
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateProcessingInvoiceForSub err:", err.Error())
			return err
		}
		createRes, err := service.CreateSubInvoicePaymentDefaultAutomatic(ctx, one, false, "", "", "SubscriptionEndTrialManual", 0)
		if err != nil {
			g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic err:", err.Error())
			return err
		}
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().CurrentPeriodStart: invoice.PeriodStart,
			dao.Subscription.Columns().CurrentPeriodEnd:   invoice.PeriodEnd,
			dao.Subscription.Columns().DunningTime:        dunningTime,
			dao.Subscription.Columns().GmtModify:          gtime.Now(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
		if err != nil {
			return err
		}
		g.Log().Print(ctx, "EndTrialManual CreateSubInvoicePaymentDefaultAutomatic:", utility.MarshalToJsonString(createRes))
		err = handler.HandleSubscriptionIncomplete(ctx, sub.SubscriptionId, gtime.Now().Timestamp())

	} else {
		_, err = dao.Subscription.Ctx(ctx).Data(g.Map{
			dao.Subscription.Columns().Status:         newStatus,
			dao.Subscription.Columns().TrialEnd:       newTrialEnd,
			dao.Subscription.Columns().DunningTime:    dunningTime,
			dao.Subscription.Columns().GmtModify:      gtime.Now(),
			dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
		}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	}
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     sub.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", sub.SubscriptionId),
		Content:        "EndTrial",
		UserId:         sub.UserId,
		SubscriptionId: sub.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		g.Log().Print(ctx, "EndTrialManual err:", err.Error())
		return err
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
		dao.Subscription.Columns().Status:         consts.SubStatusProcessing,
		dao.Subscription.Columns().GmtModify:      gtime.Now(),
		dao.Subscription.Columns().LastUpdateTime: gtime.Now().Timestamp(),
	}).Where(dao.Subscription.Columns().SubscriptionId, subscriptionId).OmitNil().Update()
	operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
		MerchantId:     one.MerchantId,
		Target:         fmt.Sprintf("Subscription(%s)", one.SubscriptionId),
		Content:        "MarkSubscriptionProcessed",
		UserId:         one.UserId,
		SubscriptionId: one.SubscriptionId,
		InvoiceId:      "",
		PlanId:         0,
		DiscountCode:   "",
	}, err)
	if err != nil {
		return err
	}
	return nil
}
