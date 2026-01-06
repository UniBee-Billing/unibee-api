package onetime

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/credit/config"
	"unibee/internal/logic/credit/payment"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/gateway/gateway_bean"
	"unibee/internal/logic/invoice/handler"
	service3 "unibee/internal/logic/invoice/service"
	handler2 "unibee/internal/logic/payment/handler"
	"unibee/internal/logic/payment/service"
	"unibee/internal/logic/plan/period"
	service2 "unibee/internal/logic/subscription/service"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/logic/user/vat"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
	"unibee/utility/unibee"
)

type SubscriptionCreateOnetimeAddonInternalReq struct {
	MerchantId             uint64                 `json:"merchantId" dc:"MerchantId"`
	SubscriptionId         string                 `json:"subscriptionId"  dc:"SubscriptionId" `
	AddonId                uint64                 `json:"addonId" dc:"addonId"`
	Currency               string                 `json:"currency"          dc:"The currency of payment"`
	Quantity               int64                  `json:"quantity" dc:"Quantity"`
	ReturnUrl              string                 `json:"returnUrl"  dc:"ReturnUrl" `
	CancelUrl              string                 `json:"cancelUrl"  dc:"CancelUrl" `
	Metadata               map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
	DiscountCode           string                 `json:"discountCode"        dc:"DiscountCode, ignore if discountAmount or discountPercentage provide"`
	DiscountAmount         *int64                 `json:"discountAmount"     dc:"Amount of discount"`
	DiscountPercentage     *int64                 `json:"discountPercentage" dc:"Percentage of discount, 100=1%, ignore if discountAmount provide"`
	TaxPercentage          *int64                 `json:"taxPercentage" dc:"TaxPercentage，1000 = 10%, use subscription's taxPercentage if not provide"`
	GatewayId              *uint64                `json:"gatewayId" dc:"GatewayId, use subscription's gateway if not provide"`
	GatewayPaymentType     string                 `json:"gatewayPaymentType" dc:"GatewayPaymentType" `
	PaymentMethodId        string                 `json:"paymentMethodId" dc:"PaymentMethodId" `
	ApplyPromoCredit       *bool                  `json:"applyPromoCredit"  dc:"apply promo credit or not"`
	ApplyPromoCreditAmount *int64                 `json:"applyPromoCreditAmount"  dc:"apply promo credit amount, auto compute if not specified"`
	IsSubmit               bool
}

type SubscriptionCreateOnetimeAddonInternalRes struct {
	MerchantId               uint64                         `json:"merchantId" dc:"MerchantId"`
	SubscriptionOnetimeAddon *bean.SubscriptionOnetimeAddon `json:"subscriptionOnetimeAddon"  dc:"SubscriptionOnetimeAddon" `
	Paid                     bool                           `json:"paid"`
	Link                     string                         `json:"link"`
	Invoice                  *bean.Invoice                  `json:"invoice"  dc:"Invoice" `
}

type SubscriptionCreateOnetimeAddonPreviewInternalRes struct {
	MerchantId           uint64                  `json:"merchantId" dc:"MerchantId"`
	Subscription         *entity.Subscription    `json:"subscription" dc:"Subscription"`
	User                 *entity.UserAccount     `json:"user" dc:"User"`
	Addon                *entity.Plan            `json:"addon" dc:"Addon"`
	Quantity             int64                   `json:"quantity" dc:"Quantity"`
	Gateway              *entity.MerchantGateway `json:"gateway" dc:"Gateway"`
	GatewayId            uint64                  `json:"gatewayId" dc:"GatewayId"`
	GatewayPaymentType   string                  `json:"gatewayPaymentType" dc:"GatewayPaymentType" `
	GatewayPaymentMethod string                  `json:"gatewayPaymentMethod" dc:"GatewayPaymentMethod" `
	Invoice              *bean.Invoice           `json:"invoice"  dc:"Invoice" `
	DiscountMessage      string                  `json:"discountMessage" `
	ApplyPromoCredit     bool                    `json:"applyPromoCredit" `
}

func CreateSubOneTimeAddonPreview(ctx context.Context, req *SubscriptionCreateOnetimeAddonInternalReq) (*SubscriptionCreateOnetimeAddonPreviewInternalRes, error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.AddonId > 0, "AddonId invalid")
	req.Quantity = utility.MaxInt64(req.Quantity, 1)
	addon := query.GetPlanById(ctx, req.AddonId)
	utility.Assert(addon != nil, "Addon not found")
	utility.Assert(addon.Status == consts.PlanStatusActive, "Addon not active")
	utility.Assert(addon.Type != consts.PlanTypeMain, "Addon not onetime type")
	currency := addon.Currency
	if len(req.Currency) > 0 {
		currency = req.Currency
	}
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "Sub not found")
	//utility.Assert(sub.Currency == addon.Currency, "Server error: currency not match")
	user := query.GetUserAccountById(ctx, sub.UserId)
	utility.Assert(user != nil, "User not found")
	gatewayId, paymentType, paymentMethodId := sub_update.VerifyPaymentGatewayMethod(ctx, sub.UserId, req.GatewayId, req.GatewayPaymentType, req.PaymentMethodId, sub.SubscriptionId)
	utility.Assert(gatewayId > 0, "Gateway need specified")
	if req.GatewayId != nil {
		gatewayId = *req.GatewayId
	}
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, "Gateway not found")
	var taxPercentage = sub.TaxPercentage
	percentage, countryCode, vatNumber, err := vat.GetUserTaxPercentage(ctx, sub.UserId)
	if err == nil {
		taxPercentage = percentage
	}
	if req.TaxPercentage != nil {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "External TaxPercentage only available for api call")
		utility.Assert(*req.TaxPercentage >= 0 && *req.TaxPercentage < 10000, "invalid taxPercentage")
		taxPercentage = *req.TaxPercentage
	}
	promoCreditDiscountCodeExclusive := config.CheckCreditConfigDiscountCodeExclusive(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, currency)

	var discountMessage string
	if len(req.DiscountCode) > 0 {
		canApply, _, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:         sub.MerchantId,
			UserId:             sub.UserId,
			DiscountCode:       req.DiscountCode,
			Currency:           currency,
			PLanId:             req.AddonId,
			TimeNow:            gtime.Now().Timestamp(),
			IsUpgrade:          false,
			IsChangeToLongPlan: false,
			IsRenew:            false,
			IsNewUser:          service2.IsNewSubscriptionUser(ctx, sub.MerchantId, user.Email),
		})
		if canApply {

		} else {
			req.DiscountCode = ""
			discountMessage = message
		}
		{
			//conflict, disable discount code
			if promoCreditDiscountCodeExclusive && canApply && req.ApplyPromoCredit != nil && *req.ApplyPromoCredit {
				_, promoCreditPayout, _ := payment.CheckCreditUserPayout(ctx, req.MerchantId, sub.UserId, consts.CreditAccountTypePromo, currency, addon.CurrencyAmount(ctx, currency), req.ApplyPromoCreditAmount)
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

	totalAmountExcludingTax := addon.CurrencyAmount(ctx, currency) * req.Quantity

	if req.ApplyPromoCredit == nil {
		if promoCreditDiscountCodeExclusive && len(req.DiscountCode) > 0 {
			req.ApplyPromoCredit = unibee.Bool(false)
		} else {
			req.ApplyPromoCredit = unibee.Bool(config.CheckCreditConfigPreviewDefaultUsed(ctx, _interface.GetMerchantId(ctx), consts.CreditAccountTypePromo, currency))
		}
	}

	//Promo Credit
	var promoCreditDiscountAmount int64 = 0
	var promoCreditAccount *bean.CreditAccount
	var promoCreditPayout *bean.CreditPayout
	var creditPayoutErr error
	if *req.ApplyPromoCredit {
		promoCreditAccount, promoCreditPayout, creditPayoutErr = payment.CheckCreditUserPayout(ctx, req.MerchantId, sub.UserId, consts.CreditAccountTypePromo, currency, totalAmountExcludingTax, req.ApplyPromoCreditAmount)
		if creditPayoutErr == nil && promoCreditAccount != nil && promoCreditPayout != nil {
			promoCreditDiscountAmount = promoCreditPayout.CurrencyAmount
			totalAmountExcludingTax = totalAmountExcludingTax - promoCreditDiscountAmount
		}
	}

	var discountAmount int64 = 0
	if req.DiscountAmount != nil && *req.DiscountAmount > 0 {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		discountAmount = utility.MinInt64(*req.DiscountAmount, totalAmountExcludingTax)
	} else if req.DiscountPercentage != nil && *req.DiscountPercentage > 0 {
		utility.Assert(_interface.Context().Get(ctx).IsOpenApiCall, "Discount only available for api call")
		utility.Assert(*req.DiscountPercentage > 0 && *req.DiscountPercentage <= 10000, "invalid discountPercentage")
		discountAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(*req.DiscountPercentage))
	} else if len(req.DiscountCode) > 0 {
		discountCode := query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode)
		utility.Assert(discountCode.Type == 0, "invalid code, code is from external")
		canApply, isRecurring, message := discount.UserDiscountApplyPreview(ctx, &discount.UserDiscountApplyReq{
			MerchantId:         req.MerchantId,
			UserId:             sub.UserId,
			DiscountCode:       req.DiscountCode,
			Currency:           currency,
			PLanId:             addon.Id,
			TimeNow:            utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock),
			IsUpgrade:          false,
			IsChangeToLongPlan: false,
			IsRenew:            false,
			IsNewUser:          service2.IsNewSubscriptionUser(ctx, req.MerchantId, user.Email),
		})
		utility.Assert(canApply, message)
		utility.Assert(!isRecurring, "recurring discount code not available for one-time addon")
		discountAmount = utility.MinInt64(discount.ComputeDiscountAmount(ctx, query.GetDiscountByCode(ctx, req.MerchantId, req.DiscountCode), totalAmountExcludingTax, currency, gtime.Now().Timestamp()), totalAmountExcludingTax)
	}

	totalAmountExcludingTax = totalAmountExcludingTax - discountAmount
	var taxAmount = int64(float64(totalAmountExcludingTax) * utility.ConvertTaxPercentageToInternalFloat(taxPercentage))
	invoice := &bean.Invoice{
		InvoiceName:                    "OneTimeAddonPurchase-Subscription",
		BizType:                        consts.BizTypeOneTime,
		SubscriptionId:                 sub.SubscriptionId,
		OriginAmount:                   totalAmountExcludingTax + taxAmount + discountAmount + promoCreditDiscountAmount,
		TotalAmount:                    totalAmountExcludingTax + taxAmount,
		DiscountCode:                   req.DiscountCode,
		DiscountAmount:                 discountAmount,
		PromoCreditDiscountAmount:      promoCreditDiscountAmount,
		PromoCreditAccount:             promoCreditAccount,
		PromoCreditPayout:              promoCreditPayout,
		TotalAmountExcludingTax:        totalAmountExcludingTax,
		SubscriptionAmount:             totalAmountExcludingTax + discountAmount + promoCreditDiscountAmount + taxAmount,
		SubscriptionAmountExcludingTax: totalAmountExcludingTax + discountAmount + promoCreditDiscountAmount,
		Currency:                       currency,
		VatNumber:                      vatNumber,
		CountryCode:                    countryCode,
		TaxPercentage:                  taxPercentage,
		TaxAmount:                      taxAmount,
		PaymentType:                    paymentType,
		PaymentMethodId:                paymentMethodId,
		Lines: []*bean.InvoiceItemSimplify{{
			Currency:               currency,
			OriginAmount:           totalAmountExcludingTax + taxAmount + discountAmount + promoCreditDiscountAmount,
			Amount:                 totalAmountExcludingTax + taxAmount,
			AmountExcludingTax:     totalAmountExcludingTax,
			DiscountAmount:         discountAmount + promoCreditDiscountAmount,
			Tax:                    taxAmount,
			UnitAmountExcludingTax: addon.CurrencyAmount(ctx, currency),
			Name:                   addon.PlanName,
			Description:            addon.Description,
			Quantity:               req.Quantity,
			Plan:                   bean.SimplifyPlan(addon),
		}},
	}
	return &SubscriptionCreateOnetimeAddonPreviewInternalRes{
		MerchantId:           req.MerchantId,
		Subscription:         sub,
		User:                 user,
		Addon:                addon,
		Quantity:             req.Quantity,
		Gateway:              gateway,
		GatewayId:            gatewayId,
		GatewayPaymentType:   paymentType,
		GatewayPaymentMethod: paymentMethodId,
		Invoice:              invoice,
		DiscountMessage:      discountMessage,
		ApplyPromoCredit:     *req.ApplyPromoCredit,
	}, nil
}

func CreateSubOneTimeAddon(ctx context.Context, req *SubscriptionCreateOnetimeAddonInternalReq) (*SubscriptionCreateOnetimeAddonInternalRes, error) {
	req.IsSubmit = true
	preview, err := CreateSubOneTimeAddonPreview(ctx, req)
	if err != nil {
		return nil, err
	}
	one := &entity.SubscriptionOnetimeAddon{
		UserId:         preview.Subscription.UserId,
		SubscriptionId: req.SubscriptionId,
		AddonId:        req.AddonId,
		Quantity:       req.Quantity,
		Status:         1,
		CreateTime:     gtime.Now().Timestamp(),
		MetaData:       utility.MarshalToJsonString(req.Metadata),
	}

	result, err := dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`CreateSubOneTimeAddon SubscriptionPendingUpdate record insert failure %s`, err)
		return nil, err
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	user := query.GetUserAccountById(ctx, preview.Subscription.UserId)
	utility.Assert(user != nil, "user not found")
	utility.Assert(user.Status != 2, "Your account has been suspended")

	invoice, err := service3.CreateProcessingInvoiceForSub(ctx, &service3.CreateProcessingInvoiceForSubReq{
		PlanId:             preview.Addon.Id,
		Simplify:           preview.Invoice,
		Sub:                preview.Subscription,
		GatewayId:          preview.Gateway.Id,
		GatewayPaymentType: preview.GatewayPaymentType,
		PaymentMethodId:    preview.GatewayPaymentMethod,
		IsSubLatestInvoice: false,
		TimeNow:            gtime.Now().Timestamp(),
	})
	utility.Assert(err == nil, fmt.Sprintf("%+v", err))
	preview.Invoice.Id = invoice.Id
	preview.Invoice.InvoiceId = invoice.InvoiceId

	var createRes *gateway_bean.GatewayNewPaymentResp
	var paymentId string
	if preview.Invoice.TotalAmount > 0 {
		createRes, err = service.GatewayPaymentCreate(ctx, &gateway_bean.GatewayNewPaymentReq{
			CheckoutMode: false,
			Gateway:      preview.Gateway,
			Pay: &entity.Payment{
				ExternalPaymentId:    strconv.FormatUint(one.Id, 10),
				BizType:              consts.BizTypeOneTime,
				SubscriptionId:       preview.Subscription.SubscriptionId,
				UserId:               preview.Subscription.UserId,
				GatewayId:            preview.Gateway.Id,
				GatewayPaymentMethod: preview.GatewayPaymentMethod,
				GatewayEdition:       preview.GatewayPaymentType,
				TotalAmount:          preview.Invoice.TotalAmount,
				Currency:             preview.Invoice.Currency,
				CountryCode:          preview.Invoice.CountryCode,
				MerchantId:           preview.MerchantId,
				CompanyId:            0,
				ReturnUrl:            req.ReturnUrl,
			},
			Email:                preview.User.Email,
			Metadata:             map[string]interface{}{"BillingReason": preview.Invoice.InvoiceName, "Source": "CreateSubOneTimeAddon", "CancelUrl": req.CancelUrl, "SubscriptionOnetimeAddonId": strconv.FormatUint(one.Id, 10)},
			Invoice:              preview.Invoice,
			PayImmediate:         true,
			GatewayPaymentType:   preview.GatewayPaymentType,
			GatewayPaymentMethod: preview.GatewayPaymentMethod,
		})
		utility.Assert(err == nil, fmt.Sprintf("%+v", err))
		paymentId = createRes.Payment.PaymentId
	} else {
		invoice, err = handler.MarkInvoiceAsPaidForZeroPayment(ctx, invoice.InvoiceId)
		utility.AssertError(err, "System Error")
		sub_update.UpdateUserCountryCode(ctx, one.UserId, invoice.CountryCode)
		err = handler2.CreateOrUpdatePaymentItemForPaymentInvoice(ctx, invoice, consts.PaymentSuccess)
		if err != nil {
			g.Log().Errorf(ctx, "CreateSubOneTimeAddon CreateOrUpdatePaymentItemForPaymentInvoice error:%s", err.Error())
		}
		createRes = &gateway_bean.GatewayNewPaymentResp{
			Status:  consts.PaymentSuccess,
			Invoice: invoice,
		}
	}
	//update paymentId
	status := 1
	if createRes.Status == consts.PaymentSuccess {
		status = 2
	}
	periodStart := gtime.Now().Timestamp()
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	if sub != nil {
		periodStart = utility.MaxInt64(gtime.Now().Timestamp(), sub.TestClock)
	}
	periodEnd := period.GetPeriodEndFromStart(ctx, periodStart, periodStart, req.AddonId)
	_, err = dao.SubscriptionOnetimeAddon.Ctx(ctx).Data(g.Map{
		dao.SubscriptionOnetimeAddon.Columns().Status:      status,
		dao.SubscriptionOnetimeAddon.Columns().PeriodStart: periodStart,
		dao.SubscriptionOnetimeAddon.Columns().PeriodEnd:   periodEnd,
		dao.SubscriptionOnetimeAddon.Columns().InvoiceId:   invoice.InvoiceId,
		dao.SubscriptionOnetimeAddon.Columns().PaymentId:   paymentId,
		dao.SubscriptionOnetimeAddon.Columns().GmtModify:   gtime.Now(),
	}).Where(dao.SubscriptionOnetimeAddon.Columns().Id, one.Id).OmitNil().Update()
	if err != nil {
		return nil, err
	}

	return &SubscriptionCreateOnetimeAddonInternalRes{
		SubscriptionOnetimeAddon: bean.SimplifySubscriptionOnetimeAddon(ctx, one),
		Link:                     createRes.Link,
		Paid:                     createRes.Status == consts.PaymentSuccess,
		Invoice:                  bean.SimplifyInvoice(createRes.Invoice),
	}, nil
}
