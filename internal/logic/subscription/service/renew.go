package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strconv"
	"unibee/api/bean"
	config2 "unibee/internal/cmd/config"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/discount"
	"unibee/internal/logic/gateway/gateway_bean"
	handler2 "unibee/internal/logic/invoice/handler"
	"unibee/internal/logic/invoice/invoice_compute"
	service3 "unibee/internal/logic/invoice/service"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/payment/service"
	subscription2 "unibee/internal/logic/subscription"
	"unibee/internal/logic/subscription/handler"
	"unibee/internal/logic/subscription/pending_update_cancel"
	"unibee/internal/logic/user/sub_update"
	"unibee/internal/query"
	"unibee/utility"
)

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
			utility.AssertError(err, "Create Gateway Payment Error")
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
		dao.Subscription.Columns().TrialEnd:  sub.CurrentPeriodStart - 1,
		dao.Subscription.Columns().GmtModify: gtime.Now(),
	}).Where(dao.Subscription.Columns().SubscriptionId, sub.SubscriptionId).OmitNil().Update()

	return &CreateInternalRes{
		Subscription: bean.SimplifySubscription(ctx, sub),
		Paid:         createRes.Status == consts.PaymentSuccess,
		Link:         createRes.Link,
	}, nil
}