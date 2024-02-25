package api

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/balance"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/invoice"
	"github.com/stripe/stripe-go/v76/invoiceitem"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/refund"
	sub "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/taxrate"
	"strconv"
	"strings"
	"time"
	"unibee-api/internal/consts"
	dao "unibee-api/internal/dao/oversea_pay"
	webhook2 "unibee-api/internal/logic/gateway"
	"unibee-api/internal/logic/gateway/api/log"
	_ "unibee-api/internal/logic/gateway/base"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/gateway/util"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type Stripe struct {
}

func (s Stripe) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := queryAndCreateChannelUser(ctx, gateway, userId)

	params := &stripe.CustomerListPaymentMethodsParams{
		Customer: stripe.String(gatewayUser.GatewayUserId),
	}
	params.Limit = stripe.Int64(10)
	result := customer.ListPaymentMethods(params)
	var paymentMethods = make([]string, 0)
	for _, paymentMethod := range result.PaymentMethodList().Data {
		paymentMethods = append(paymentMethods, paymentMethod.ID)
	}
	return &ro.GatewayUserPaymentMethodListInternalResp{
		PaymentMethods: paymentMethods,
	}, nil
}

func (s Stripe) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.CustomerParams{
		//Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
		Email: stripe.String(user.Email),
	}

	createCustomResult, err := customer.New(params)
	log.SaveChannelHttpLog("GatewayUserCreate", params, createCustomResult, err, "", nil, gateway)
	if err != nil {
		g.Log().Printf(ctx, "customer.New: %v", err.Error())
		return nil, err
	}
	return &ro.GatewayUserCreateInternalResp{GatewayUserId: createCustomResult.ID}, nil

}

func (s Stripe) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := queryAndCreateChannelUser(ctx, gateway, listReq.UserId)

	params := &stripe.PaymentIntentListParams{}
	params.Customer = stripe.String(gatewayUser.GatewayUserId)
	params.Limit = stripe.Int64(200)
	paymentList := paymentintent.List(params)
	log.SaveChannelHttpLog("GatewayPaymentList", params, paymentList, err, "", nil, gateway)
	var list []*ro.GatewayPaymentRo
	for _, item := range paymentList.PaymentIntentList().Data {
		list = append(list, parseStripePayment(item, gateway))
	}

	return list, nil
}

func (s Stripe) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.RefundListParams{}
	params.PaymentIntent = stripe.String(gatewayPaymentId)
	params.Limit = stripe.Int64(100)
	refundList := refund.List(params)
	log.SaveChannelHttpLog("GatewayRefundList", params, refundList, err, "", nil, gateway)
	var list []*ro.OutPayRefundRo
	for _, item := range refundList.RefundList().Data {
		list = append(list, parseStripeRefund(item))
	}

	return list, nil
}

func (s Stripe) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.PaymentIntentParams{}
	response, err := paymentintent.Get(gatewayPaymentId, params)
	log.SaveChannelHttpLog("GatewayPaymentDetail", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}

	return parseStripePayment(response, gateway), nil
}

func (s Stripe) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{}
	response, err := refund.Get(gatewayRefundId, params)
	log.SaveChannelHttpLog("GatewayRefundDetail", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return parseStripeRefund(response), nil
}

func (s Stripe) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.BalanceParams{}
	response, err := balance.Get(params)
	if err != nil {
		return nil, err
	}

	var availableBalances []*ro.GatewayBalance
	for _, item := range response.Available {
		availableBalances = append(availableBalances, &ro.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var connectReservedBalances []*ro.GatewayBalance
	for _, item := range response.ConnectReserved {
		connectReservedBalances = append(connectReservedBalances, &ro.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var pendingBalances []*ro.GatewayBalance
	for _, item := range response.ConnectReserved {
		pendingBalances = append(pendingBalances, &ro.GatewayBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	return &ro.GatewayMerchantBalanceQueryInternalResp{
		AvailableBalance:       availableBalances,
		ConnectReservedBalance: connectReservedBalances,
		PendingBalance:         pendingBalances,
	}, nil
}

func (s Stripe) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.CustomerParams{}
	response, err := customer.Get(queryAndCreateChannelUserWithOutPaymentMethod(ctx, gateway, userId).GatewayUserId, params)
	if err != nil {
		return nil, err
	}
	var cashBalances []*ro.GatewayBalance
	if response.CashBalance != nil {
		for currency, amount := range response.CashBalance.Available {
			cashBalances = append(cashBalances, &ro.GatewayBalance{
				Amount:   amount,
				Currency: strings.ToUpper(currency),
			})
		}
	}

	var invoiceCreditBalances []*ro.GatewayBalance
	for currency, amount := range response.InvoiceCreditBalance {
		invoiceCreditBalances = append(invoiceCreditBalances, &ro.GatewayBalance{
			Amount:   amount,
			Currency: strings.ToUpper(currency),
		})
	}
	var defaultPaymentMethod string
	if response.InvoiceSettings != nil && response.InvoiceSettings.DefaultPaymentMethod != nil {
		defaultPaymentMethod = response.InvoiceSettings.DefaultPaymentMethod.ID
	}
	return &ro.GatewayUserDetailQueryInternalResp{
		GatewayUserId:        response.ID,
		DefaultPaymentMethod: defaultPaymentMethod,
		Balance: &ro.GatewayBalance{
			Amount:   response.Balance,
			Currency: strings.ToUpper(string(response.Currency)),
		},
		CashBalance:          cashBalances,
		InvoiceCreditBalance: invoiceCreditBalances,
		Description:          response.Description,
		Email:                response.Email,
	}, nil
}

func (s Stripe) GatewaySubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		TrialEndNow:       stripe.Bool(true),
		ProrationBehavior: stripe.String("none"),
	}
	_, err = sub.Update(subscription.GatewaySubscriptionId, params)
	if err != nil {
		return nil, err
	}

	details, err := s.GatewaySubscriptionDetails(ctx, plan, gatewayPlan, subscription)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// GatewaySubscriptionNewTrialEnd https://stripe.com/docs/billing/subscriptions/billing-cycle#add-a-trial-to-change-the-billing-cycle
func (s Stripe) GatewaySubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		//TrialEnd:          stripe.Int64(newTrialEnd),
		BillingCycleAnchor: stripe.Int64(newTrialEnd), // todo mark test for use anchor
		ProrationBehavior:  stripe.String("none"),
	}
	_, err = sub.Update(subscription.GatewaySubscriptionId, params)
	if err != nil {
		return nil, err
	}

	details, err := s.GatewaySubscriptionDetails(ctx, plan, gatewayPlan, subscription)
	if err != nil {
		return nil, err
	}
	if details.TrialEnd != newTrialEnd {
		return nil, gerror.New("update new trial end error")
	}
	return details, nil
}

// Test Card Data
// Payment Success
// 4242 4242 4242 4242
// Need 3DS
// 4000 0025 0000 3155
// Reject
// 4000 0000 0000 9995
func (s Stripe) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func (s Stripe) GatewaySubscriptionCreate(ctx context.Context, subscriptionRo *ro.GatewayCreateSubscriptionInternalReq) (res *ro.GatewayCreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.GatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, subscriptionRo.GatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	{
		gatewayUser := queryAndCreateChannelUser(ctx, gateway, subscriptionRo.Subscription.UserId)

		gatewayVatRate := query.GetSubscriptionVatRateChannel(ctx, subscriptionRo.VatCountryRate.Id, gateway.Id)
		if gatewayVatRate == nil {
			params := &stripe.TaxRateParams{
				DisplayName: stripe.String("VAT"),
				Description: stripe.String(subscriptionRo.VatCountryRate.CountryName),
				Percentage:  stripe.Float64(utility.ConvertTaxScaleToPercentageFloat(subscriptionRo.VatCountryRate.StandardTaxPercentage)),
				Country:     stripe.String(subscriptionRo.VatCountryRate.CountryCode),
				Active:      stripe.Bool(true),
				//Jurisdiction: stripe.String("DE"),
				Inclusive: stripe.Bool(false),
			}
			vatCreateResult, err := taxrate.New(params)
			if err != nil {
				g.Log().Printf(ctx, "taxrate.New: %v", err.Error())
				return nil, err
			}
			gatewayVatRate = &entity.GatewayVatRate{
				VatRateId:        int64(subscriptionRo.VatCountryRate.Id),
				GatewayId:        int64(gateway.Id),
				GatewayVatRateId: vatCreateResult.ID,
				CreateTime:       gtime.Now().Timestamp(),
			}
			result, err := dao.GatewayVatRate.Ctx(ctx).Data(gatewayVatRate).OmitNil().Insert(gatewayVatRate)
			if err != nil {
				err = gerror.Newf(`SubscriptionVatRateChannel record insert failure %s`, err.Error())
				return nil, err
			}
			id, _ := result.LastInsertId()
			gatewayVatRate.Id = uint64(uint(id))
		}

		var checkoutMode = true
		if checkoutMode {
			items := []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.CheckoutSessionLineItemParams{
					Price:    stripe.String(addon.AddonGatewayPlan.GatewayPlanId),
					Quantity: stripe.Int64(addon.Quantity),
				})
			}
			checkoutParams := &stripe.CheckoutSessionParams{
				Customer:  stripe.String(gatewayUser.GatewayUserId),
				Currency:  stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
				LineItems: items,
				//AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
				//	Enabled: stripe.Bool(!taxInclusive), //Default值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_rates
				//},
				Metadata: map[string]string{
					"SubId": subscriptionRo.Subscription.SubscriptionId,
				},
				SuccessURL: stripe.String(webhook2.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, true)),
				CancelURL:  stripe.String(webhook2.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, false)),
			}
			checkoutParams.Mode = stripe.String(string(stripe.CheckoutSessionModeSubscription))
			checkoutParams.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
				Metadata: map[string]string{
					"SubId": subscriptionRo.Subscription.SubscriptionId,
				},
				DefaultTaxRates: []*string{stripe.String(gatewayVatRate.GatewayVatRateId)},
			}
			//checkoutParams.ExpiresAt
			createSubscription, err := session.New(checkoutParams)
			log.SaveChannelHttpLog("GatewaySubscriptionCreateSession", checkoutParams, createSubscription, err, "", nil, gateway)
			if err != nil {
				return nil, err
			}
			return &ro.GatewayCreateSubscriptionInternalResp{
				GatewayUserId: gatewayUser.GatewayUserId,
				Link:          createSubscription.URL,
				Data:          utility.FormatToJsonString(createSubscription),
				Status:        0, //todo mark
			}, nil
		} else {
			items := []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Main",
						"BillingPlanId":   strconv.FormatUint(subscriptionRo.GatewayPlan.PlanId, 10),
					},
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.SubscriptionItemsParams{
					Price:    stripe.String(addon.AddonGatewayPlan.GatewayPlanId),
					Quantity: stripe.Int64(addon.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Addon",
						"BillingPlanId":   strconv.FormatUint(addon.AddonGatewayPlan.PlanId, 10),
					},
				})
			}
			subscriptionParams := &stripe.SubscriptionParams{
				Customer: stripe.String(gatewayUser.GatewayUserId),
				Currency: stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
				Items:    items,
				//AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
				//	Enabled: stripe.Bool(!taxInclusive), //Default false = need stripe compute tax，true != todo
				//},
				PaymentBehavior:  stripe.String("default_incomplete"),   // todo mark https://stripe.com/docs/api/subscriptions/create
				CollectionMethod: stripe.String("charge_automatically"), //Default charge_automatically，charge automatic
				Metadata: map[string]string{
					"SubId": subscriptionRo.Subscription.SubscriptionId,
				},
				DefaultTaxRates: []*string{stripe.String(gatewayVatRate.GatewayVatRateId)},
			}
			subscriptionParams.AddExpand("latest_invoice.payment_intent")
			createSubscription, err := sub.New(subscriptionParams)
			log.SaveChannelHttpLog("GatewaySubscriptionCreate", subscriptionParams, createSubscription, err, "", nil, gateway)
			if err != nil {
				return nil, err
			}

			return &ro.GatewayCreateSubscriptionInternalResp{
				GatewayUserId:             gatewayUser.GatewayUserId,
				Link:                      createSubscription.LatestInvoice.HostedInvoiceURL,
				GatewaySubscriptionId:     createSubscription.ID,
				GatewaySubscriptionStatus: string(createSubscription.Status),
				Data:                      utility.FormatToJsonString(createSubscription),
				Status:                    0, //todo mark
				Paid:                      createSubscription.LatestInvoice.Paid,
			}, nil
		}
	}
	//{
	//	//付款链接方式可能存在多次重复付款问题
	//	params := &stripe.PaymentLinkParams{
	//		LineItems: []*stripe.PaymentLinkLineItemParams{
	//			{
	//				Price:    stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
	//				Quantity: stripe.Int64(1),
	//			},
	//			//{
	//			//	Price: stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
	//			//},
	//		},
	//		AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
	//			Type: stripe.String(string(stripe.PaymentLinkAfterCompletionTypeRedirect)),
	//			Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
	//				URL: stripe.String("https://www.baidu.com"),
	//			},
	//		},
	//		//不启用试用期
	//		//SubscriptionData: &stripe.PaymentLinkSubscriptionDataParams{
	//		//	TrialPeriodDays: stripe.Int64(7),
	//		//},
	//	}
	//	createSubscription, err := paymentlink.New(params)
	//	if err != nil {
	//		return nil, err
	//	}
	//	jsonData, _ := gjson.Marshal(createSubscription)
	//	return &ro.GatewayCreateSubscriptionInternalResp{
	//		GatewaySubscriptionId:     createSubscription.ID,
	//		GatewaySubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}
	//{
	//	checkoutParams := &stripe.CheckoutSessionParams{
	//		Metadata: map[string]string{
	//			"orderId": subscriptionRo.Subscription.SubscriptionId,
	//		},
	//		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	//		LineItems: []*stripe.CheckoutSessionLineItemParams{
	//			{
	//				Price:    stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
	//				Quantity: stripe.Int64(1),
	//			},
	//		},
	//		SuccessURL: stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, true)),
	//		CancelURL:  stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, false)),
	//	}
	//
	//	result, err := session.New(checkoutParams)
	//	if err != nil {
	//		return nil, err
	//	}
	//	jsonData, _ := gjson.Marshal(result)
	//	return &ro.GatewayCreateSubscriptionInternalResp{
	//		GatewaySubscriptionId:     result.ID,
	//		GatewaySubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}

}

// GatewaySubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel?dashboard-or-api=api
func (s Stripe) GatewaySubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.GatewayCancelSubscriptionInternalReq) (res *ro.GatewayCancelSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionCancelInternalReq.Subscription.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, subscriptionCancelInternalReq.Subscription.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionCancelParams{}
	params.InvoiceNow = stripe.Bool(subscriptionCancelInternalReq.InvoiceNow)
	params.Prorate = stripe.Bool(subscriptionCancelInternalReq.Prorate)
	response, err := sub.Cancel(subscriptionCancelInternalReq.Subscription.GatewaySubscriptionId, params)
	log.SaveChannelHttpLog("GatewaySubscriptionCancel", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return &ro.GatewayCancelSubscriptionInternalResp{}, nil
}

// GatewaySubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel
func (s Stripe) GatewaySubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.GatewaySubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.GatewaySubscriptionId, params)
	log.SaveChannelHttpLog("GatewaySubscriptionCancelAtPeriodEnd", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return &ro.GatewayCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

func (s Stripe) GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.GatewaySubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.GatewaySubscriptionId, params)
	log.SaveChannelHttpLog("GatewaySubscriptionCancelLastCancelAtPeriodEnd", params, response, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return &ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

func (s Stripe) GatewaySubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionPreviewInternalResp, err error) {
	utility.Assert(subscriptionRo.GatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, subscriptionRo.GatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	gatewayUser := queryAndCreateChannelUser(ctx, gateway, subscriptionRo.Subscription.UserId)

	// Set the proration date to this moment:
	updateUnixTime := time.Now().Unix()
	if consts.NonEffectImmediatelyUsePendingUpdate && !subscriptionRo.EffectImmediate {
		updateUnixTime = subscriptionRo.Subscription.CurrentPeriodEnd
	}
	if subscriptionRo.ProrationDate > 0 {
		updateUnixTime = subscriptionRo.ProrationDate
	}
	items, err := s.makeSubscriptionUpdateItems(subscriptionRo)
	if err != nil {
		return nil, err
	}
	params := &stripe.InvoiceUpcomingParams{
		Customer:          stripe.String(gatewayUser.GatewayUserId),
		Subscription:      stripe.String(subscriptionRo.Subscription.GatewaySubscriptionId),
		SubscriptionItems: items,
		//SubscriptionProrationBehavior: stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice)),// 设置了就只会输出 Proration 账单
	}
	params.SubscriptionProrationDate = stripe.Int64(updateUnixTime)
	detail, err := invoice.Upcoming(params)
	log.SaveChannelHttpLog("GatewaySubscriptionUpdateProrationPreview", params, detail, err, subscriptionRo.Subscription.GatewaySubscriptionId, nil, gateway)
	if err != nil {
		return nil, err
	}

	// 拆开 invoice Proration into invoice,nextPeriodInvoice
	var currentInvoiceItems []*ro.InvoiceItemDetailRo
	var nextInvoiceItems []*ro.InvoiceItemDetailRo
	var currentSubAmount int64 = 0
	var currentSubAmountExcludingTax int64 = 0
	var nextSubAmount int64 = 0
	var nextSubAmountExcludingTax int64 = 0
	for _, line := range detail.Lines.Data {
		if line.Proration {
			currentInvoiceItems = append(currentInvoiceItems, &ro.InvoiceItemDetailRo{
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
				Description:            line.Description,
				Proration:              line.Proration,
				Quantity:               line.Quantity,
				Currency:               strings.ToUpper(string(line.Currency)),
			})
			currentSubAmount = currentSubAmount + line.Amount
			currentSubAmountExcludingTax = currentSubAmountExcludingTax + line.AmountExcludingTax
		} else {
			nextInvoiceItems = append(nextInvoiceItems, &ro.InvoiceItemDetailRo{
				Amount:                 line.Amount,
				AmountExcludingTax:     line.AmountExcludingTax,
				UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
				Description:            line.Description,
				Proration:              line.Proration,
				Quantity:               line.Quantity,
				Currency:               strings.ToUpper(string(line.Currency)),
			})
			nextSubAmount = nextSubAmount + line.Amount
			nextSubAmountExcludingTax = nextSubAmountExcludingTax + line.AmountExcludingTax
		}
	}

	currentInvoice := &ro.GatewayDetailInvoiceInternalResp{
		TotalAmount:                    currentSubAmount,
		TotalAmountExcludingTax:        currentSubAmountExcludingTax,
		TaxAmount:                      currentSubAmount - currentSubAmountExcludingTax,
		SubscriptionAmount:             currentSubAmount,
		SubscriptionAmountExcludingTax: currentSubAmountExcludingTax,
		Lines:                          currentInvoiceItems,
		GatewaySubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		GatewayId:                      int64(gateway.Id),
		GatewayUserId:                  detail.Customer.ID,
	}

	nextPeriodInvoice := &ro.GatewayDetailInvoiceInternalResp{
		TotalAmount:                    nextSubAmount,
		TotalAmountExcludingTax:        nextSubAmountExcludingTax,
		TaxAmount:                      nextSubAmount - nextSubAmountExcludingTax,
		SubscriptionAmount:             nextSubAmount,
		SubscriptionAmountExcludingTax: nextSubAmountExcludingTax,
		Lines:                          nextInvoiceItems,
		GatewaySubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		GatewayId:                      int64(gateway.Id),
		GatewayUserId:                  detail.Customer.ID,
	}

	return &ro.GatewayUpdateSubscriptionPreviewInternalResp{
		Data:          utility.FormatToJsonString(detail),
		TotalAmount:   currentInvoice.TotalAmount,
		Currency:      strings.ToUpper(string(detail.Currency)),
		ProrationDate: updateUnixTime,
		Invoice:       currentInvoice,
		//Invoice: parseStripeInvoice(detail, int64(gateway.Id)),
		NextPeriodInvoice: nextPeriodInvoice,
	}, nil
}

func (s Stripe) makeSubscriptionUpdateItems(subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) ([]*stripe.SubscriptionItemsParams, error) {

	var items []*stripe.SubscriptionItemsParams

	var stripeSubscriptionItems []*stripe.SubscriptionItem
	if !subscriptionRo.EffectImmediate && !consts.NonEffectImmediatelyUsePendingUpdate {
		if len(subscriptionRo.Subscription.GatewayItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.GatewayItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.GatewaySubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//Solution 1 range and delete，Effect Next Period，PendingUpdate Not Support
		for _, item := range stripeSubscriptionItems {
			//delete all，Add Plan and Addons
			items = append(items, &stripe.SubscriptionItemsParams{
				ID:      stripe.String(item.ID),
				Deleted: stripe.Bool(true),
			})
		}
		//Add Plan
		items = append(items, &stripe.SubscriptionItemsParams{
			Price:    stripe.String(subscriptionRo.GatewayPlan.GatewayPlanId),
			Quantity: stripe.Int64(subscriptionRo.Quantity),
			Metadata: map[string]string{
				"BillingPlanType": "Main",
				"BillingPlanId":   strconv.FormatUint(subscriptionRo.GatewayPlan.PlanId, 10),
			},
		})
		for _, addon := range subscriptionRo.AddonPlans {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(addon.AddonGatewayPlan.GatewayPlanId),
				Quantity: stripe.Int64(addon.Quantity),
				Metadata: map[string]string{
					"BillingPlanType": "Addon",
					"BillingPlanId":   strconv.FormatUint(addon.AddonGatewayPlan.PlanId, 10),
				},
			})
		}
	} else {
		//Use PendingUpdate
		if len(subscriptionRo.Subscription.GatewayItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.GatewayItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.GatewaySubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//Solution 2 EffectImmediate=true, Use PendingUpdate，Modify Quantity = 0 for Plan&Addon Need Delete，
		newMap := make(map[string]int64)
		for _, addon := range subscriptionRo.AddonPlans {
			newMap[addon.AddonGatewayPlan.GatewayPlanId] = addon.Quantity
		}
		newMap[subscriptionRo.GatewayPlan.GatewayPlanId] = subscriptionRo.Quantity
		//Range Match
		for _, item := range stripeSubscriptionItems {
			if quantity, ok := newMap[item.Price.ID]; ok {
				//Replace
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(item.Price.ID),
					Quantity: stripe.Int64(quantity),
				})
				delete(newMap, item.Price.ID)
			} else {
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Quantity: stripe.Int64(0),
				})
			}
		}
		//Add Others
		for GatewayPlanId, quantity := range newMap {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(GatewayPlanId),
				Quantity: stripe.Int64(quantity),
			})
		}
	}

	return items, nil
}

// GatewaySubscriptionUpdate Price Can Not Duplicate In Items
func (s Stripe) GatewaySubscriptionUpdate(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.GatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, subscriptionRo.GatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	items, err := s.makeSubscriptionUpdateItems(subscriptionRo)
	if err != nil {
		return nil, err
	}

	params := &stripe.SubscriptionParams{
		Items: items,
	}
	if subscriptionRo.EffectImmediate {
		if consts.ProrationUsingUniBeeCompute {
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorNone))
		} else {
			params.ProrationDate = stripe.Int64(subscriptionRo.ProrationDate)
			params.PaymentBehavior = stripe.String("pending_if_incomplete") //pendingIfIncomplete Only Some Attr Can Update,  Price Quantity
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice))
		}
	} else {
		if consts.NonEffectImmediatelyUsePendingUpdate {
			params.ProrationDate = stripe.Int64(subscriptionRo.ProrationDate)
			params.PaymentBehavior = stripe.String("pending_if_incomplete") //pendingIfIncomplete Only Some Attr Can Update,  Price Quantity
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice))
		} else {
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorNone))
		}
	}
	updateSubscription, err := sub.Update(subscriptionRo.Subscription.GatewaySubscriptionId, params)
	log.SaveChannelHttpLog("GatewaySubscriptionUpdate", params, updateSubscription, err, subscriptionRo.Subscription.GatewaySubscriptionId, nil, gateway)
	if err != nil {
		return nil, err
	}

	if subscriptionRo.EffectImmediate && !consts.ProrationUsingUniBeeCompute {
		queryParams := &stripe.InvoiceParams{}
		newInvoice, err := invoice.Get(updateSubscription.LatestInvoice.ID, queryParams)
		log.SaveChannelHttpLog("GatewaySubscriptionUpdate", queryParams, newInvoice, err, "GetInvoice", nil, gateway)
		g.Log().Infof(ctx, "query new invoice:%v", newInvoice)

		return &ro.GatewayUpdateSubscriptionInternalResp{
			Data:            utility.FormatToJsonString(updateSubscription),
			GatewayUpdateId: newInvoice.ID,
			Link:            newInvoice.HostedInvoiceURL,
			Paid:            newInvoice.Paid,
		}, nil
	} else {
		//EffectImmediate=false Do Not Need Pay, The Invoice Is Old
		return &ro.GatewayUpdateSubscriptionInternalResp{
			Data: utility.FormatToJsonString(updateSubscription),
			Paid: false,
		}, nil
	}
}

// GatewaySubscriptionDetails，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
func (s Stripe) GatewaySubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.SubscriptionParams{}
	response, err := sub.Get(subscription.GatewaySubscriptionId, params)
	log.SaveChannelHttpLog("GatewaySubscriptionDetails", params, response, err, subscription.GatewaySubscriptionId, nil, gateway)
	if err != nil {
		return nil, err
	}
	//var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	//if strings.Compare(string(response.Status), "trialing") == 0 ||
	//	strings.Compare(string(response.Status), "active") == 0 {
	//	status = consts.SubStatusActive
	//} else if strings.Compare(string(response.Status), "incomplete") == 0 ||
	//	strings.Compare(string(response.Status), "incomplete_expired") == 0 {
	//	status = consts.SubStatusCreate
	//} else if strings.Compare(string(response.Status), "past_due") == 0 ||
	//	strings.Compare(string(response.Status), "unpaid") == 0 ||
	//	strings.Compare(string(response.Status), "paused") == 0 {
	//	status = consts.SubStatusSuspended
	//} else if strings.Compare(string(response.Status), "canceled") == 0 {
	//	status = consts.SubStatusCancelled
	//}

	return parseStripeSubscription(response), nil
}

// GatewayPlanActive 使用 price 代替 plan  https://stripe.com/docs/api/plans
func (s Stripe) GatewayPlanActive(ctx context.Context, targetPlan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(true) // todo mark 使用这种方式可能不能用
	result, err := price.Update(gatewayPlan.GatewayPlanId, params)
	log.SaveChannelHttpLog("GatewayPlanActive", params, result, err, "", nil, gateway)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) GatewayPlanDeactivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(false) // todo mark this way may not work
	result, err := price.Update(gatewayPlan.GatewayPlanId, params)
	log.SaveChannelHttpLog("GatewayPlanDeactivate", params, result, err, "", nil, gateway)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) GatewayProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreateProductInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.ProductParams{
		Active:      stripe.Bool(true),
		Description: stripe.String(plan.GatewayProductDescription), // todo mark not sure about if description is nil
		Name:        stripe.String(plan.GatewayProductName),
	}
	if len(plan.ImageUrl) > 0 {
		params.Images = stripe.StringSlice([]string{plan.ImageUrl})
	}
	if len(plan.HomeUrl) > 0 {
		params.URL = stripe.String(plan.HomeUrl)
	}
	result, err := product.New(params)
	log.SaveChannelHttpLog("GatewayProductCreate", params, result, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	//Prod Status Seems Not Active After Create todo mark
	return &ro.GatewayCreateProductInternalResp{
		GatewayProductId:     result.ID,
		GatewayProductStatus: fmt.Sprintf("%v", result.Active),
	}, nil
}

func (s Stripe) GatewayPlanCreateAndActivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreatePlanInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	//params := &stripe.PlanParams{
	//	//todo mark
	//	Active:   stripe.Bool(true),
	//	Amount:   stripe.Int64(1200),
	//	Currency: stripe.String(string(stripe.CurrencyUSD)),
	//	Interval: stripe.String(string(stripe.PlanIntervalMonth)),
	//	Product:  &stripe.PlanProductParams{ID: stripe.String("prod_NjpI7DbZx6AlWQ")},
	//}
	//result, err := plan.New(params)
	// Price Replace Plan https://stripe.com/docs/api/plans
	params := &stripe.PriceParams{
		Currency:   stripe.String(strings.ToLower(targetPlan.Currency)),
		UnitAmount: stripe.Int64(targetPlan.Amount),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(targetPlan.IntervalUnit),
			IntervalCount: stripe.Int64(int64(targetPlan.IntervalCount)),
		},
		Product: stripe.String(gatewayPlan.GatewayProductId),
		Metadata: map[string]string{
			"PlanId": strconv.FormatUint(targetPlan.Id, 10),
			"Type":   strconv.Itoa(targetPlan.Type),
		},
		//ProductData: &stripe.PriceProductDataParams{
		//	ID:   stripe.String(gatewayPlan.GatewayProductId),
		//	Name: stripe.String(targetPlan.PlanName),
		//},//this is create
	}
	result, err := price.New(params)
	log.SaveChannelHttpLog("GatewayPlanCreateAndActivate", params, result, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	return &ro.GatewayCreatePlanInternalResp{
		GatewayPlanId:     result.ID,
		GatewayPlanStatus: fmt.Sprintf("%v", result.Active),
		Data:              utility.FormatToJsonString(result),
		Status:            consts.GatewayPlanStatusActive,
	}, nil
}

func (s Stripe) GatewayInvoiceCancel(ctx context.Context, gateway *entity.MerchantGateway, cancelInvoiceInternalReq *ro.GatewayCancelInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoiceMarkUncollectibleParams{}
	response, err := invoice.MarkUncollectible(cancelInvoiceInternalReq.GatewayInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("GatewayInvoiceCancel", params, response, err, "", nil, gateway)
	return parseStripeInvoice(response, int64(gateway.Id)), nil
}

func (s Stripe) GatewayInvoiceCreateAndPay(ctx context.Context, gateway *entity.MerchantGateway, createInvoiceInternalReq *ro.GatewayCreateInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	gatewayUser := queryAndCreateChannelUser(ctx, gateway, createInvoiceInternalReq.Invoice.UserId)

	params := &stripe.InvoiceParams{
		Currency: stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
		Customer: stripe.String(gatewayUser.GatewayUserId)}
	if createInvoiceInternalReq.PayMethod == 1 {
		params.CollectionMethod = stripe.String("charge_automatically")
	} else {
		params.CollectionMethod = stripe.String("send_invoice")
		if createInvoiceInternalReq.DaysUtilDue > 0 {
			params.DaysUntilDue = stripe.Int64(int64(createInvoiceInternalReq.DaysUtilDue))
		}
	}
	//params.DefaultTaxRates
	result, err := invoice.New(params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("GatewayInvoiceCreateAndPay", params, result, err, "New", nil, gateway)

	for _, line := range createInvoiceInternalReq.InvoiceLines {
		ItemParams := &stripe.InvoiceItemParams{
			Invoice:  stripe.String(result.ID),
			Currency: stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
			//UnitAmount:  stripe.Int64(line.UnitAmountExcludingTax),
			Amount:      stripe.Int64(line.Amount),
			Description: stripe.String(line.Description),
			//Quantity:    stripe.Int64(line.Quantity),
			Customer: stripe.String(gatewayUser.GatewayUserId)}
		_, err = invoiceitem.New(ItemParams)
		if err != nil {
			return nil, err
		}
	}

	finalizeInvoiceParam := &stripe.InvoiceFinalizeInvoiceParams{}
	if createInvoiceInternalReq.PayMethod == 1 {
		finalizeInvoiceParam.AutoAdvance = stripe.Bool(true)
	} else {
		finalizeInvoiceParam.AutoAdvance = stripe.Bool(false)
	}

	detail, err := invoice.FinalizeInvoice(result.ID, finalizeInvoiceParam)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("GatewayInvoiceCreateAndPay", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, gateway)
	var status consts.InvoiceStatusEnum = consts.InvoiceStatusInit
	if strings.Compare(string(detail.Status), "draft") == 0 {
		status = consts.InvoiceStatusPending
	} else if strings.Compare(string(detail.Status), "open") == 0 {
		status = consts.InvoiceStatusProcessing
	} else if strings.Compare(string(detail.Status), "paid") == 0 {
		status = consts.InvoiceStatusPaid
	} else if strings.Compare(string(detail.Status), "uncollectible") == 0 {
		status = consts.InvoiceStatusFailed
	} else if strings.Compare(string(detail.Status), "void") == 0 {
		status = consts.InvoiceStatusCancelled
	}
	var invoiceItems []*ro.InvoiceItemDetailRo
	for _, line := range detail.Lines.Data {
		var start int64 = 0
		var end int64 = 0
		if line.Period != nil {
			start = line.Period.Start
			end = line.Period.End
		}
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               strings.ToUpper(string(line.Currency)),
			Amount:                 line.Amount,
			AmountExcludingTax:     line.AmountExcludingTax,
			UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
			Description:            line.Description,
			Proration:              line.Proration,
			Quantity:               line.Quantity,
			PeriodStart:            start,
			PeriodEnd:              end,
		})
	}

	return &ro.GatewayDetailInvoiceInternalResp{
		TotalAmount:                    detail.Total,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		GatewayId:                      int64(gateway.Id),
		Status:                         status,
		GatewayUserId:                  detail.Customer.ID,
		Link:                           detail.HostedInvoiceURL,
		GatewayStatus:                  string(detail.Status),
		GatewayInvoiceId:               detail.ID,
		GatewayInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
	}, nil
}

func (s Stripe) GatewayInvoicePay(ctx context.Context, gateway *entity.MerchantGateway, payInvoiceInternalReq *ro.GatewayPayInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoicePayParams{}
	response, err := invoice.Pay(payInvoiceInternalReq.GatewayInvoiceId, params)
	log.SaveChannelHttpLog("GatewayInvoicePay", params, response, err, "", nil, gateway)
	return parseStripeInvoice(response, int64(gateway.Id)), nil
}

func (s Stripe) GatewayInvoiceDetails(ctx context.Context, gateway *entity.MerchantGateway, gatewayInvoiceId string) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	params := &stripe.InvoiceParams{}
	detail, err := invoice.Get(gatewayInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("GatewayInvoiceDetails", params, detail, err, "", nil, gateway)
	return parseStripeInvoice(detail, int64(gateway.Id)), nil
}

func (s Stripe) GatewayPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.Gateway != nil, "gateway not found")
	stripe.Key = createPayContext.Gateway.GatewaySecret
	s.setUnibeeAppInfo()
	gatewayUser := queryAndCreateChannelUser(ctx, createPayContext.Gateway, createPayContext.Pay.UserId)

	if createPayContext.CheckoutMode {
		var items []*stripe.CheckoutSessionLineItemParams
		for _, line := range createPayContext.Invoice.Lines {
			items = append(items, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(strings.ToLower(line.Currency)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("%s", line.Description)),
					},
					UnitAmount: stripe.Int64(line.Amount),
				},
				Quantity: stripe.Int64(1),
				//TaxRates: []*string{stripe.String(gatewayVatRate.ChannelVatRateId)}, // not apply tax
			})
		}
		checkoutParams := &stripe.CheckoutSessionParams{
			Customer:   stripe.String(gatewayUser.GatewayUserId),
			Currency:   stripe.String(strings.ToLower(createPayContext.Pay.Currency)),
			LineItems:  items,
			SuccessURL: stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true)),
			CancelURL:  stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false)),
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				Metadata: createPayContext.MediaData,
			},
		}
		checkoutParams.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
		checkoutParams.Metadata = createPayContext.MediaData
		//checkoutParams.ExpiresAt
		detail, err := session.New(checkoutParams)
		if err != nil {
			return nil, err
		}
		log.SaveChannelHttpLog("GatewayPayment", checkoutParams, detail, err, "CheckoutSession", nil, createPayContext.Gateway)
		var status consts.PayStatusEnum = consts.TO_BE_PAID
		if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusOpen)) == 0 {
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusComplete)) == 0 {
			status = consts.PAY_SUCCESS
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusExpired)) == 0 {
			status = consts.PAY_FAILED
		}
		var gatewayPaymentId string
		if detail.PaymentIntent != nil {
			gatewayPaymentId = detail.PaymentIntent.ID
		}
		return &ro.CreatePayInternalResp{
			Status:                 status,
			GatewayPaymentId:       gatewayPaymentId,
			GatewayPaymentIntentId: detail.ID,
			Link:                   detail.URL,
		}, nil
	} else {
		if createPayContext.PayMethod == 1 && createPayContext.PayImmediate {
			// try use payment intent
			listQuery, err := s.GatewayUserPaymentMethodListQuery(ctx, createPayContext.Gateway, gatewayUser.UserId)
			log.SaveChannelHttpLog("GatewayPayment", gatewayUser.UserId, listQuery, err, "GatewayUserPaymentMethodListQuery", nil, createPayContext.Gateway)
			if err != nil {
				return nil, err
			}

			var success = false
			var targetIntent *stripe.PaymentIntent
			var gatewayPaymentId = ""
			var link = ""
			var cancelErr error
			if len(gatewayUser.GatewayDefaultPaymentMethod) > 0 {
				if listQuery.PaymentMethods == nil {
					listQuery.PaymentMethods = make([]string, 0)
				}
				listQuery.PaymentMethods = append(listQuery.PaymentMethods, gatewayUser.GatewayDefaultPaymentMethod)
			}
			for _, method := range listQuery.PaymentMethods {
				params := &stripe.PaymentIntentParams{
					Customer: stripe.String(gatewayUser.GatewayUserId),
					Confirm:  stripe.Bool(true),
					Amount:   stripe.Int64(createPayContext.Invoice.TotalAmount),
					Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
					AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
						Enabled: stripe.Bool(true),
					},
					Metadata:  createPayContext.MediaData,
					ReturnURL: stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true)),
				}
				params.PaymentMethod = stripe.String(method)
				targetIntent, err = paymentintent.New(params)
				log.SaveChannelHttpLog("GatewayPayment", params, targetIntent, err, "PaymentIntentCreate", nil, createPayContext.Gateway)
				var status = ""

				if targetIntent != nil {
					status = string(targetIntent.Status)
					gatewayPaymentId = targetIntent.ID
					if targetIntent.Invoice != nil {
						link = targetIntent.Invoice.HostedInvoiceURL
					}
				}
				if err == nil && strings.Compare(status, "succeeded") == 0 {
					success = true
					break
				} else if targetIntent != nil {
					cancelRes, cancelErr := paymentintent.Cancel(targetIntent.ID, &stripe.PaymentIntentCancelParams{})
					log.SaveChannelHttpLog("GatewayPayment", "", cancelRes, cancelErr, "PaymentIntentCancel", nil, createPayContext.Gateway)
				} else {
					cancelErr = gerror.Newf("targetIntent is nil")
				}
				g.Log().Printf(ctx, "GatewayPayment try PaymentIntent Method::%s gatewayPaymentId:%s status:%s success:%v link:%s error:%s cancelErr:%s\n", method, gatewayPaymentId, status, success, link, err, cancelErr)
			}
			if success && targetIntent != nil && len(gatewayPaymentId) > 0 {
				return &ro.CreatePayInternalResp{
					Status:                 consts.PAY_SUCCESS,
					GatewayPaymentId:       gatewayPaymentId,
					GatewayPaymentIntentId: targetIntent.ID,
					Link:                   link,
				}, nil
			}
		}
		// need payment link
		params := &stripe.InvoiceParams{
			Metadata: createPayContext.MediaData,
			Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
			Customer: stripe.String(gatewayUser.GatewayUserId)}

		if createPayContext.PayMethod == 1 {
			params.CollectionMethod = stripe.String("charge_automatically")
			// check the gateway user contains the payment method now
			listQuery, err := s.GatewayUserPaymentMethodListQuery(ctx, createPayContext.Gateway, gatewayUser.UserId)
			if err != nil {
				return nil, err
			}
			if len(createPayContext.GatewayPaymentMethod) > 0 && ContainString(listQuery.PaymentMethods, createPayContext.GatewayPaymentMethod) {
				params.DefaultPaymentMethod = stripe.String(createPayContext.GatewayPaymentMethod)
			} else if len(gatewayUser.GatewayDefaultPaymentMethod) > 0 && ContainString(listQuery.PaymentMethods, gatewayUser.GatewayDefaultPaymentMethod) {
				params.DefaultPaymentMethod = stripe.String(gatewayUser.GatewayDefaultPaymentMethod)
			} else if len(listQuery.PaymentMethods) > 0 {
				// todo mark use detail query
				params.DefaultPaymentMethod = stripe.String(listQuery.PaymentMethods[0])
			}
		} else {
			params.CollectionMethod = stripe.String("send_invoice")
			if createPayContext.DaysUtilDue > 0 {
				params.DaysUntilDue = stripe.Int64(int64(createPayContext.DaysUtilDue))
			}
		}
		result, err := invoice.New(params)
		if err != nil {
			return nil, err
		}
		log.SaveChannelHttpLog("GatewayPayment", params, result, err, "NewInvoice", nil, createPayContext.Gateway)

		for _, line := range createPayContext.Invoice.Lines {
			ItemParams := &stripe.InvoiceItemParams{
				Invoice:  stripe.String(result.ID),
				Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
				//UnitAmount:  stripe.Int64(line.UnitAmountExcludingTax),
				Amount:      stripe.Int64(line.Amount),
				Description: stripe.String(line.Description),
				//Quantity:    stripe.Int64(line.Quantity),
				Customer: stripe.String(gatewayUser.GatewayUserId)}
			_, err = invoiceitem.New(ItemParams)
			if err != nil {
				return nil, err
			}
		}
		finalizeInvoiceParam := &stripe.InvoiceFinalizeInvoiceParams{}
		if createPayContext.PayMethod == 1 {
			finalizeInvoiceParam.AutoAdvance = stripe.Bool(true)
		} else {
			finalizeInvoiceParam.AutoAdvance = stripe.Bool(false)
		}
		detail, err := invoice.FinalizeInvoice(result.ID, finalizeInvoiceParam)
		log.SaveChannelHttpLog("GatewayPayment", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, createPayContext.Gateway)
		if err != nil {
			return nil, err
		}
		if createPayContext.PayImmediate {
			params := &stripe.InvoicePayParams{}
			response, payErr := invoice.Pay(result.ID, params)
			log.SaveChannelHttpLog("GatewayPayment", params, response, payErr, "PayInvoice", nil, createPayContext.Gateway)
			if response != nil {
				detail.Status = response.Status
			}
		}
		var status consts.PayStatusEnum = consts.TO_BE_PAID
		if strings.Compare(string(detail.Status), "draft") == 0 {
		} else if strings.Compare(string(detail.Status), "open") == 0 {
		} else if strings.Compare(string(detail.Status), "paid") == 0 {
			status = consts.PAY_SUCCESS
		} else if strings.Compare(string(detail.Status), "uncollectible") == 0 {
			status = consts.PAY_FAILED
		} else if strings.Compare(string(detail.Status), "void") == 0 {
			status = consts.PAY_FAILED
		}
		var gatewayPaymentId string
		if detail.PaymentIntent != nil {
			gatewayPaymentId = detail.PaymentIntent.ID
		}
		return &ro.CreatePayInternalResp{
			Status:                 status,
			GatewayPaymentId:       gatewayPaymentId,
			GatewayPaymentIntentId: detail.ID,
			Link:                   detail.HostedInvoiceURL,
		}, nil
	}
}

func ContainString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (s Stripe) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	utility.Assert(payment.GatewayId > 0, "invalid payment gatewayId")
	utility.Assert(len(payment.GatewayPaymentIntentId) > 0, "invalid payment GatewayPaymentIntentId")
	gateway := util.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()

	var status = consts.TO_BE_PAID
	var gatewayCancelId string
	if strings.HasPrefix(payment.GatewayPaymentIntentId, "in_") {
		params := &stripe.InvoiceVoidInvoiceParams{}
		result, err := invoice.VoidInvoice(payment.GatewayPaymentIntentId, params)
		log.SaveChannelHttpLog("GatewayCancel", params, result, err, "VoidInvoice", nil, gateway)
		if err != nil {
			return nil, err
		}
		invoiceDetails := parseStripeInvoice(result, payment.GatewayId)
		gatewayCancelId = result.ID
		if invoiceDetails.Status == consts.InvoiceStatusPaid {
			status = consts.PAY_SUCCESS
		} else if invoiceDetails.Status == consts.InvoiceStatusFailed || invoiceDetails.Status == consts.InvoiceStatusCancelled {
			status = consts.PAY_FAILED
		}
	} else if strings.HasPrefix(payment.GatewayPaymentIntentId, "cs_") {
		params := &stripe.CheckoutSessionExpireParams{}
		result, err := session.Expire(
			payment.GatewayPaymentIntentId,
			params,
		)
		log.SaveChannelHttpLog("GatewayCancel", params, result, err, "ExpireSession", nil, gateway)
		if err != nil {
			return nil, err
		}
		if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusOpen)) == 0 {
		} else if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusComplete)) == 0 {
			status = consts.PAY_SUCCESS
		} else if strings.Compare(string(result.Status), string(stripe.CheckoutSessionStatusExpired)) == 0 {
			status = consts.PAY_FAILED
		}
		gatewayCancelId = result.ID
	} else {
		params := &stripe.PaymentIntentCancelParams{}
		result, err := paymentintent.Cancel(payment.GatewayPaymentIntentId, params)
		log.SaveChannelHttpLog("GatewayCancel", params, result, err, "CancelPaymentIntent", nil, gateway)
		if err != nil {
			return nil, err
		}
		paymentDetails := parseStripePayment(result, gateway)
		status = paymentDetails.Status
		gatewayCancelId = paymentDetails.GatewayPaymentId
	}

	return &ro.OutPayCancelRo{
		MerchantId:      strconv.FormatUint(payment.MerchantId, 10),
		GatewayCancelId: gatewayCancelId,
		Reference:       payment.PaymentId,
		Status:          strconv.Itoa(status),
	}, nil
}

func (s Stripe) GatewayPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) GatewayRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) GatewayRefund(ctx context.Context, payment *entity.Payment, one *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(payment.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, payment.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	stripe.Key = gateway.GatewaySecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{PaymentIntent: stripe.String(payment.GatewayPaymentId)}
	params.Reason = stripe.String("requested_by_customer")
	params.Amount = stripe.Int64(one.RefundAmount)
	//params.Currency = stripe.String(strings.ToLower(one.Currency))
	params.Metadata = map[string]string{"RefundId": one.RefundId}
	result, err := refund.New(params)
	log.SaveChannelHttpLog("GatewayRefund", params, result, err, "refund", nil, gateway)
	utility.Assert(err == nil, fmt.Sprintf("call stripe refund error %s", err))
	utility.Assert(result != nil, "Stripe refund failed, result is nil")
	return &ro.OutPayRefundRo{
		GatewayRefundId: result.ID,
		Status:          consts.REFUND_ING,
	}, nil
}

func parseStripeRefund(item *stripe.Refund) *ro.OutPayRefundRo {
	var gatewayPaymentId string
	if item.PaymentIntent != nil {
		gatewayPaymentId = item.PaymentIntent.ID
	}
	var status consts.RefundStatusEnum = consts.REFUND_ING
	if strings.Compare(string(item.Status), "succeeded") == 0 {
		status = consts.REFUND_SUCCESS
	} else if strings.Compare(string(item.Status), "failed") == 0 {
		status = consts.REFUND_FAILED
	} else if strings.Compare(string(item.Status), "canceled") == 0 {
		status = consts.REFUND_REVERSE
	}
	return &ro.OutPayRefundRo{
		MerchantId:       "",
		GatewayRefundId:  item.ID,
		GatewayPaymentId: gatewayPaymentId,
		Status:           status,
		Reason:           string(item.Reason),
		RefundAmount:     item.Amount,
		Currency:         strings.ToUpper(string(item.Currency)),
		RefundTime:       gtime.NewFromTimeStamp(item.Created),
	}
}

func parseStripePayment(item *stripe.PaymentIntent, gateway *entity.MerchantGateway) *ro.GatewayPaymentRo {
	var gatewayInvoiceId string
	if item.Invoice != nil {
		gatewayInvoiceId = item.Invoice.ID
	}
	var gatewayUserId string
	if item.Customer != nil {
		gatewayUserId = item.Customer.ID
	}
	var status = consts.TO_BE_PAID
	if strings.Compare(string(item.Status), "succeeded") == 0 {
		status = consts.PAY_SUCCESS
	} else if strings.Compare(string(item.Status), "canceled") == 0 {
		status = consts.PAY_CANCEL
	}
	var captureStatus = consts.AUTHORIZED
	var authorizeReason = ""
	var paymentData = ""
	if strings.Compare(string(item.Status), "requires_payment_method") == 0 {
		captureStatus = consts.WAITING_AUTHORIZED
		if item.LastPaymentError != nil {
			authorizeReason = item.LastPaymentError.Msg
		}
	} else if strings.Compare(string(item.Status), "requires_confirmation") == 0 {
		captureStatus = consts.CAPTURE_REQUEST
	}
	if item.NextAction != nil {
		paymentData = utility.MarshalToJsonString(item.NextAction)
	}
	var gatewayPaymentMethod string
	if item.PaymentMethod != nil {
		gatewayPaymentMethod = item.PaymentMethod.ID
	}
	return &ro.GatewayPaymentRo{
		GatewayId:            int64(gateway.Id),
		MerchantId:           gateway.MerchantId,
		GatewayInvoiceId:     gatewayInvoiceId,
		GatewayUserId:        gatewayUserId,
		GatewayPaymentId:     item.ID,
		Status:               status,
		AuthorizeStatus:      captureStatus,
		AuthorizeReason:      authorizeReason,
		CancelReason:         string(item.CancellationReason),
		PaymentData:          paymentData,
		TotalAmount:          item.Amount,
		PaymentAmount:        item.AmountReceived,
		GatewayPaymentMethod: gatewayPaymentMethod,
		Currency:             strings.ToUpper(string(item.Currency)),
		PayTime:              gtime.NewFromTimeStamp(item.Created),
		CreateTime:           gtime.NewFromTimeStamp(item.Created),
		CancelTime:           gtime.NewFromTimeStamp(item.CanceledAt),
	}
}

func parseStripeSubscription(subscription *stripe.Subscription) *ro.GatewayDetailSubscriptionInternalResp {
	//https://stripe.com/docs/billing/subscriptions/overview
	/**
	trialing	订阅目前处于试用期，可以安全地为您的客户配置您的产品。订阅会自动转换到active首次付款时。
	active	订阅信誉良好，最近一次付款成功。为您的客户配置您的产品是安全的。
	incomplete	需要在23小时内成功付款才能激活订阅。或者付款需要采取行动，例如客户身份验证。incomplete如果有待付款并且 PaymentIntent 状态为 ，则订阅也可以为processing。
	incomplete_expired	订阅的首次付款失败，并且在创建订阅后 23 小时内未成功付款。这些订阅不会向客户收取费用。存在此状态是为了让您可以跟踪未能激活订阅的客户。
	past_due	最新最终发票的付款失败或未尝试。订阅将继续创建发票。您的订阅设置决定了订阅的下一个状态。如果在尝试所有智能重试后发票仍未支付，您可以将订阅配置为移至canceled、unpaid，或保留为past_due。要将订阅转移到active，请在到期日之前支付最新的发票。
	canceled	订阅已被取消。取消期间，将禁用所有未付发票的自动收取 ( auto_advance=false)。这是无法更新的最终状态。
	unpaid	最新的发票尚未支付，但订阅仍然有效。最新发票仍处于打开状态，并且继续生成发票，但不会尝试付款。您应该在订阅时撤销对产品的访问权限，unpaid因为已尝试付款并在订阅时重试past_due。要将订阅转移到active，请在到期日之前支付最新的发票。
	paused	订阅已结束试用期，没有默认付款方式，并且trial_settings.end_behavior.missing_payment_method设置为pause。将不再为订阅创建发票。为客户附加默认付款方式后，您可以恢复订阅。
	*/
	var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	if strings.Compare(string(subscription.Status), "trialing") == 0 ||
		strings.Compare(string(subscription.Status), "active") == 0 {
		status = consts.SubStatusActive
	} else if strings.Compare(string(subscription.Status), "unpaid") == 0 {
		status = consts.SubStatusCreate
	} else if strings.Compare(string(subscription.Status), "incomplete_expired") == 0 {
		status = consts.SubStatusExpired
	} else if strings.Compare(string(subscription.Status), "incomplete") == 0 ||
		strings.Compare(string(subscription.Status), "pass_due") == 0 {
		status = consts.SubStatusIncomplete
	} else if strings.Compare(string(subscription.Status), "paused") == 0 {
		status = consts.SubStatusSuspended
	} else if strings.Compare(string(subscription.Status), "canceled") == 0 {
		status = consts.SubStatusCancelled
	}
	var latestChannelPaymentId = ""
	if subscription.LatestInvoice != nil && subscription.LatestInvoice.PaymentIntent != nil {
		latestChannelPaymentId = subscription.LatestInvoice.PaymentIntent.ID
	}
	var gatewayDefaultPaymentMethod = ""
	if subscription.DefaultPaymentMethod != nil {
		gatewayDefaultPaymentMethod = subscription.DefaultPaymentMethod.ID
	}

	return &ro.GatewayDetailSubscriptionInternalResp{
		Status:                      status,
		GatewaySubscriptionId:       subscription.ID,
		GatewayStatus:               string(subscription.Status),
		Data:                        utility.FormatToJsonString(subscription),
		GatewayItemData:             utility.MarshalToJsonString(subscription.Items.Data),
		GatewayLatestInvoiceId:      subscription.LatestInvoice.ID,
		GatewayLatestPaymentId:      latestChannelPaymentId,
		GatewayDefaultPaymentMethod: gatewayDefaultPaymentMethod,
		CancelAtPeriodEnd:           subscription.CancelAtPeriodEnd,
		CurrentPeriodStart:          subscription.CurrentPeriodStart,
		CurrentPeriodEnd:            subscription.CurrentPeriodEnd,
		BillingCycleAnchor:          subscription.BillingCycleAnchor,
		TrialEnd:                    subscription.TrialEnd,
	}
}

func parseStripeInvoice(detail *stripe.Invoice, gatewayId int64) *ro.GatewayDetailInvoiceInternalResp {
	var status consts.InvoiceStatusEnum = consts.InvoiceStatusInit
	if strings.Compare(string(detail.Status), "draft") == 0 {
		status = consts.InvoiceStatusPending
	} else if strings.Compare(string(detail.Status), "open") == 0 {
		status = consts.InvoiceStatusProcessing
	} else if strings.Compare(string(detail.Status), "paid") == 0 {
		status = consts.InvoiceStatusPaid
	} else if strings.Compare(string(detail.Status), "uncollectible") == 0 {
		status = consts.InvoiceStatusFailed
	} else if strings.Compare(string(detail.Status), "void") == 0 {
		status = consts.InvoiceStatusCancelled
	}
	var invoiceItems []*ro.InvoiceItemDetailRo
	for _, line := range detail.Lines.Data {
		var start int64 = 0
		var end int64 = 0
		if line.Period != nil {
			start = line.Period.Start
			end = line.Period.End
		}
		invoiceItems = append(invoiceItems, &ro.InvoiceItemDetailRo{
			Currency:               strings.ToUpper(string(line.Currency)),
			Amount:                 line.Amount,
			AmountExcludingTax:     line.AmountExcludingTax,
			UnitAmountExcludingTax: int64(line.UnitAmountExcludingTax),
			Description:            line.Description,
			Proration:              line.Proration,
			Quantity:               line.Quantity,
			PeriodStart:            start,
			PeriodEnd:              end,
		})
	}

	var gatewayPaymentId string
	if detail.PaymentIntent != nil {
		gatewayPaymentId = detail.PaymentIntent.ID
	}
	var gatewaySubscriptionId string
	if detail.Subscription != nil {
		gatewaySubscriptionId = detail.Subscription.ID
	}
	var subscriptionId string
	if detail.SubscriptionDetails != nil {
		subscriptionId = detail.SubscriptionDetails.Metadata["SubId"]
	}
	var gatewayUserId string
	if detail.Customer != nil {
		gatewayUserId = detail.Customer.ID
	}
	var paymentTime int64
	var cancelTime int64
	if detail.StatusTransitions != nil {
		paymentTime = detail.StatusTransitions.PaidAt
		cancelTime = detail.StatusTransitions.VoidedAt
	}
	var gatewayDefaultPaymentMethod = ""
	if detail.DefaultPaymentMethod != nil {
		gatewayDefaultPaymentMethod = detail.DefaultPaymentMethod.ID
	}

	return &ro.GatewayDetailInvoiceInternalResp{
		GatewayDefaultPaymentMethod:    gatewayDefaultPaymentMethod,
		TotalAmount:                    detail.Total,
		PaymentAmount:                  detail.AmountPaid,
		BalanceAmount:                  -(detail.StartingBalance) - -(detail.EndingBalance),
		BalanceStart:                   -detail.StartingBalance,
		BalanceEnd:                     -detail.EndingBalance,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		GatewayId:                      gatewayId,
		Status:                         status,
		Link:                           detail.HostedInvoiceURL,
		GatewayStatus:                  string(detail.Status),
		GatewayInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
		GatewayInvoiceId:               detail.ID,
		GatewayUserId:                  gatewayUserId,
		GatewaySubscriptionId:          gatewaySubscriptionId,
		SubscriptionId:                 subscriptionId,
		GatewayPaymentId:               gatewayPaymentId,
		PaymentTime:                    paymentTime,
		Reason:                         string(detail.BillingReason),
		CreateTime:                     detail.Created,
		CancelTime:                     cancelTime,
	}
}
