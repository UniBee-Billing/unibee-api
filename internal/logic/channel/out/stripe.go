package out

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
	"go-oversea-pay/internal/consts"
	dao "go-oversea-pay/internal/dao/oversea_pay"
	webhook2 "go-oversea-pay/internal/logic/channel"
	_ "go-oversea-pay/internal/logic/channel/base"
	"go-oversea-pay/internal/logic/channel/out/log"
	"go-oversea-pay/internal/logic/channel/ro"
	"go-oversea-pay/internal/logic/channel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
	"time"
)

type Stripe struct {
}

func (s Stripe) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	channelUser := queryAndCreateChannelUser(ctx, payChannel, userId)

	params := &stripe.CustomerListPaymentMethodsParams{
		Customer: stripe.String(channelUser.ChannelUserId),
	}
	params.Limit = stripe.Int64(10)
	result := customer.ListPaymentMethods(params)
	var paymentMethods = make([]string, 0)
	for _, paymentMethod := range result.PaymentMethodList().Data {
		paymentMethods = append(paymentMethods, paymentMethod.ID)
	}
	return &ro.ChannelUserPaymentMethodListInternalResp{
		PaymentMethods: paymentMethods,
	}, nil
}

func (s Stripe) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantChannelConfig, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.CustomerParams{
		//Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
		Email: stripe.String(user.Email),
	}

	createCustomResult, err := customer.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelUserCreate", params, createCustomResult, err, "", nil, payChannel)
	if err != nil {
		g.Log().Printf(ctx, "customer.New: %v", err.Error())
		return nil, err
	}
	return &ro.ChannelUserCreateInternalResp{ChannelUserId: createCustomResult.ID}, nil

}

func (s Stripe) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantChannelConfig, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	channelUser := queryAndCreateChannelUser(ctx, payChannel, listReq.UserId)

	params := &stripe.PaymentIntentListParams{}
	params.Customer = stripe.String(channelUser.ChannelUserId)
	params.Limit = stripe.Int64(200)
	paymentList := paymentintent.List(params)
	log.SaveChannelHttpLog("DoRemoteChannelPaymentList", params, paymentList, err, "", nil, payChannel)
	var list []*ro.ChannelPaymentRo
	for _, item := range paymentList.PaymentIntentList().Data {
		list = append(list, parseStripePayment(item, payChannel))
	}

	return list, nil
}

func (s Stripe) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.RefundListParams{}
	params.PaymentIntent = stripe.String(channelPaymentId)
	params.Limit = stripe.Int64(100)
	refundList := refund.List(params)
	log.SaveChannelHttpLog("DoRemoteChannelRefundList", params, refundList, err, "", nil, payChannel)
	var list []*ro.OutPayRefundRo
	for _, item := range refundList.RefundList().Data {
		list = append(list, parseStripeRefund(item))
	}

	return list, nil
}

func (s Stripe) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PaymentIntentParams{}
	response, err := paymentintent.Get(channelPaymentId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPaymentDetail", params, response, err, "", nil, payChannel)
	if err != nil {
		return nil, err
	}

	return parseStripePayment(response, payChannel), nil
}

func (s Stripe) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.RefundParams{}
	response, err := refund.Get(channelRefundId, params)
	log.SaveChannelHttpLog("DoRemoteChannelRefundDetail", params, response, err, "", nil, payChannel)
	if err != nil {
		return nil, err
	}
	return parseStripeRefund(response), nil
}

func (s Stripe) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.BalanceParams{}
	response, err := balance.Get(params)
	if err != nil {
		return nil, err
	}

	var availableBalances []*ro.ChannelBalance
	for _, item := range response.Available {
		availableBalances = append(availableBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var connectReservedBalances []*ro.ChannelBalance
	for _, item := range response.ConnectReserved {
		connectReservedBalances = append(connectReservedBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	var pendingBalances []*ro.ChannelBalance
	for _, item := range response.ConnectReserved {
		pendingBalances = append(pendingBalances, &ro.ChannelBalance{
			Amount:   item.Amount,
			Currency: strings.ToUpper(string(item.Currency)),
		})
	}
	return &ro.ChannelMerchantBalanceQueryInternalResp{
		AvailableBalance:       availableBalances,
		ConnectReservedBalance: connectReservedBalances,
		PendingBalance:         pendingBalances,
	}, nil
}

func (s Stripe) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantChannelConfig, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.CustomerParams{}
	response, err := customer.Get(queryAndCreateChannelUserWithOutPaymentMethod(ctx, payChannel, userId).ChannelUserId, params)
	if err != nil {
		return nil, err
	}
	var cashBalances []*ro.ChannelBalance
	if response.CashBalance != nil {
		for currency, amount := range response.CashBalance.Available {
			cashBalances = append(cashBalances, &ro.ChannelBalance{
				Amount:   amount,
				Currency: strings.ToUpper(currency),
			})
		}
	}

	var invoiceCreditBalances []*ro.ChannelBalance
	for currency, amount := range response.InvoiceCreditBalance {
		invoiceCreditBalances = append(invoiceCreditBalances, &ro.ChannelBalance{
			Amount:   amount,
			Currency: strings.ToUpper(currency),
		})
	}
	var defaultPaymentMethod string
	if response.InvoiceSettings != nil && response.InvoiceSettings.DefaultPaymentMethod != nil {
		defaultPaymentMethod = response.InvoiceSettings.DefaultPaymentMethod.ID
	}
	return &ro.ChannelUserDetailQueryInternalResp{
		ChannelUserId:        response.ID,
		DefaultPaymentMethod: defaultPaymentMethod,
		Balance: &ro.ChannelBalance{
			Amount:   response.Balance,
			Currency: strings.ToUpper(string(response.Currency)),
		},
		CashBalance:          cashBalances,
		InvoiceCreditBalance: invoiceCreditBalances,
		Description:          response.Description,
		Email:                response.Email,
	}, nil
}

func (s Stripe) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		TrialEndNow:       stripe.Bool(true),
		ProrationBehavior: stripe.String("none"),
	}
	_, err = sub.Update(subscription.ChannelSubscriptionId, params)
	if err != nil {
		return nil, err
	}

	details, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, subscription)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// DoRemoteChannelSubscriptionNewTrialEnd https://stripe.com/docs/billing/subscriptions/billing-cycle#add-a-trial-to-change-the-billing-cycle
func (s Stripe) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionParams{
		//TrialEnd:          stripe.Int64(newTrialEnd),
		BillingCycleAnchor: stripe.Int64(newTrialEnd), // todo mark test for use anchor
		ProrationBehavior:  stripe.String("none"),
	}
	_, err = sub.Update(subscription.ChannelSubscriptionId, params)
	if err != nil {
		return nil, err
	}

	details, err := s.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, subscription)
	if err != nil {
		return nil, err
	}
	if details.TrialEnd != newTrialEnd {
		return nil, gerror.New("update new trial end error")
	}
	return details, nil
}

// 测试数据
// 付款成功
// 4242 4242 4242 4242
// 付款需要验证
// 4000 0025 0000 3155
// 付款被拒绝
// 4000 0000 0000 9995
func (s Stripe) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func (s Stripe) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	{
		channelUser := queryAndCreateChannelUser(ctx, channelEntity, subscriptionRo.Subscription.UserId)

		channelVatRate := query.GetSubscriptionVatRateChannel(ctx, subscriptionRo.VatCountryRate.Id, channelEntity.Id)
		if channelVatRate == nil {
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
			channelVatRate = &entity.ChannelVatRate{
				VatRateId:        int64(subscriptionRo.VatCountryRate.Id),
				ChannelId:        int64(channelEntity.Id),
				ChannelVatRateId: vatCreateResult.ID,
			}
			result, err := dao.ChannelVatRate.Ctx(ctx).Data(channelVatRate).OmitNil().Insert(channelVatRate)
			if err != nil {
				err = gerror.Newf(`SubscriptionVatRateChannel record insert failure %s`, err.Error())
				return nil, err
			}
			id, _ := result.LastInsertId()
			channelVatRate.Id = uint64(uint(id))
		}

		//taxInclusive := true
		//if subscriptionRo.Plan.TaxInclusive == 0 {
		//	//税费不包含
		//	taxInclusive = false
		//}

		var checkoutMode = true
		if checkoutMode {
			items := []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.CheckoutSessionLineItemParams{
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
				})
			}
			checkoutParams := &stripe.CheckoutSessionParams{
				Customer:  stripe.String(channelUser.ChannelUserId),
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
				DefaultTaxRates: []*string{stripe.String(channelVatRate.ChannelVatRateId)},
			}
			//checkoutParams.ExpiresAt
			createSubscription, err := session.New(checkoutParams)
			log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCreateSession", checkoutParams, createSubscription, err, "", nil, channelEntity)
			if err != nil {
				return nil, err
			}
			return &ro.ChannelCreateSubscriptionInternalResp{
				ChannelUserId: channelUser.ChannelUserId,
				Link:          createSubscription.URL,
				Data:          utility.FormatToJsonString(createSubscription),
				Status:        0, //todo mark
			}, nil
		} else {
			items := []*stripe.SubscriptionItemsParams{
				{
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(subscriptionRo.Subscription.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Main",
						"BillingPlanId":   strconv.FormatInt(subscriptionRo.PlanChannel.PlanId, 10),
					},
				},
			}
			for _, addon := range subscriptionRo.AddonPlans {
				items = append(items, &stripe.SubscriptionItemsParams{
					Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(addon.Quantity),
					Metadata: map[string]string{
						"BillingPlanType": "Addon",
						"BillingPlanId":   strconv.FormatInt(addon.AddonPlanChannel.PlanId, 10),
					},
				})
			}
			subscriptionParams := &stripe.SubscriptionParams{
				Customer: stripe.String(channelUser.ChannelUserId),
				Currency: stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
				Items:    items,
				//AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
				//	Enabled: stripe.Bool(!taxInclusive), //Default值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_tates
				//},
				PaymentBehavior:  stripe.String("default_incomplete"),   // todo mark https://stripe.com/docs/api/subscriptions/create
				CollectionMethod: stripe.String("charge_automatically"), //Default行为 charge_automatically，自动扣款
				Metadata: map[string]string{
					"SubId": subscriptionRo.Subscription.SubscriptionId,
				},
				DefaultTaxRates: []*string{stripe.String(channelVatRate.ChannelVatRateId)},
			}
			subscriptionParams.AddExpand("latest_invoice.payment_intent")
			createSubscription, err := sub.New(subscriptionParams)
			log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCreate", subscriptionParams, createSubscription, err, "", nil, channelEntity)
			if err != nil {
				return nil, err
			}

			return &ro.ChannelCreateSubscriptionInternalResp{
				ChannelUserId:             channelUser.ChannelUserId,
				Link:                      createSubscription.LatestInvoice.HostedInvoiceURL,
				ChannelSubscriptionId:     createSubscription.ID,
				ChannelSubscriptionStatus: string(createSubscription.Status),
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
	//				Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
	//				Quantity: stripe.Int64(1),
	//			},
	//			//{
	//			//	Price: stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
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
	//	return &ro.ChannelCreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     createSubscription.ID,
	//		ChannelSubscriptionStatus: "true",
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
	//				Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
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
	//	return &ro.ChannelCreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     result.ID,
	//		ChannelSubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}

}

// DoRemoteChannelSubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel?dashboard-or-api=api
func (s Stripe) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionCancelInternalReq.Subscription.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionCancelInternalReq.Subscription.ChannelId)
	utility.Assert(channelEntity != nil, "out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.SubscriptionCancelParams{}
	params.InvoiceNow = stripe.Bool(subscriptionCancelInternalReq.InvoiceNow)
	params.Prorate = stripe.Bool(subscriptionCancelInternalReq.Prorate)
	response, err := sub.Cancel(subscriptionCancelInternalReq.Subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancel", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelSubscriptionInternalResp{}, nil
}

// DoRemoteChannelSubscriptionCancel https://stripe.com/docs/billing/subscriptions/cancel
func (s Stripe) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.ChannelSubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancelAtPeriodEnd", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

func (s Stripe) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//params := &stripe.SubscriptionCancelParams{}
	//response, err := sub.Cancel(subscription.ChannelSubscriptionId, params)
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(false)} //使用更新方式取代取消接口
	response, err := sub.Update(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd", params, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp{}, nil
}

func (s Stripe) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "out channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	channelUser := queryAndCreateChannelUser(ctx, channelEntity, subscriptionRo.Subscription.UserId)

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
		Customer:          stripe.String(channelUser.ChannelUserId),
		Subscription:      stripe.String(subscriptionRo.Subscription.ChannelSubscriptionId),
		SubscriptionItems: items,
		//SubscriptionProrationBehavior: stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice)),// 设置了就只会输出 Proration 账单
	}
	params.SubscriptionProrationDate = stripe.Int64(updateUnixTime)
	detail, err := invoice.Upcoming(params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdateProrationPreview", params, detail, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
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

	currentInvoice := &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    currentSubAmount,
		TotalAmountExcludingTax:        currentSubAmountExcludingTax,
		TaxAmount:                      currentSubAmount - currentSubAmountExcludingTax,
		SubscriptionAmount:             currentSubAmount,
		SubscriptionAmountExcludingTax: currentSubAmountExcludingTax,
		Lines:                          currentInvoiceItems,
		ChannelSubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		ChannelId:                      int64(channelEntity.Id),
		ChannelUserId:                  detail.Customer.ID,
	}

	nextPeriodInvoice := &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    nextSubAmount,
		TotalAmountExcludingTax:        nextSubAmountExcludingTax,
		TaxAmount:                      nextSubAmount - nextSubAmountExcludingTax,
		SubscriptionAmount:             nextSubAmount,
		SubscriptionAmountExcludingTax: nextSubAmountExcludingTax,
		Lines:                          nextInvoiceItems,
		ChannelSubscriptionId:          detail.Subscription.ID,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		ChannelId:                      int64(channelEntity.Id),
		ChannelUserId:                  detail.Customer.ID,
	}

	return &ro.ChannelUpdateSubscriptionPreviewInternalResp{
		Data:          utility.FormatToJsonString(detail),
		TotalAmount:   currentInvoice.TotalAmount,
		Currency:      strings.ToUpper(string(detail.Currency)),
		ProrationDate: updateUnixTime,
		Invoice:       currentInvoice,
		//Invoice: parseStripeInvoice(detail, int64(channelEntity.Id)),
		NextPeriodInvoice: nextPeriodInvoice,
	}, nil
}

func (s Stripe) makeSubscriptionUpdateItems(subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) ([]*stripe.SubscriptionItemsParams, error) {

	var items []*stripe.SubscriptionItemsParams

	var stripeSubscriptionItems []*stripe.SubscriptionItem
	if !subscriptionRo.EffectImmediate && !consts.NonEffectImmediatelyUsePendingUpdate {
		if len(subscriptionRo.Subscription.ChannelItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.ChannelItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.ChannelSubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//方案 1 遍历并删除，下周期生效，不支持 PendingUpdate
		for _, item := range stripeSubscriptionItems {
			//删除之前全部，新增 Plan 和 Addons 方式
			items = append(items, &stripe.SubscriptionItemsParams{
				ID:      stripe.String(item.ID),
				Deleted: stripe.Bool(true),
			})
		}
		//新增新的项目
		items = append(items, &stripe.SubscriptionItemsParams{
			Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
			Quantity: stripe.Int64(subscriptionRo.Quantity),
			Metadata: map[string]string{
				"BillingPlanType": "Main",
				"BillingPlanId":   strconv.FormatInt(subscriptionRo.PlanChannel.PlanId, 10),
			},
		})
		for _, addon := range subscriptionRo.AddonPlans {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(addon.AddonPlanChannel.ChannelPlanId),
				Quantity: stripe.Int64(addon.Quantity),
				Metadata: map[string]string{
					"BillingPlanType": "Addon",
					"BillingPlanId":   strconv.FormatInt(addon.AddonPlanChannel.PlanId, 10),
				},
			})
		}
	} else {
		//使用PendingUpdate
		if len(subscriptionRo.Subscription.ChannelItemData) > 0 {
			err := utility.UnmarshalFromJsonString(subscriptionRo.Subscription.ChannelItemData, &stripeSubscriptionItems)
			if err != nil {
				return nil, err
			}
		} else {
			detail, err := sub.Get(subscriptionRo.Subscription.ChannelSubscriptionId, &stripe.SubscriptionParams{})
			if err != nil {
				return nil, err
			}
			stripeSubscriptionItems = detail.Items.Data
		}
		//方案 2 EffectImmediate=true, 使用PendingUpdate，对于删除的 Plan 和 Addon，修改 Quantity 为 0
		newMap := make(map[string]int64)
		for _, addon := range subscriptionRo.AddonPlans {
			newMap[addon.AddonPlanChannel.ChannelPlanId] = addon.Quantity
		}
		newMap[subscriptionRo.PlanChannel.ChannelPlanId] = subscriptionRo.Quantity
		//匹配
		for _, item := range stripeSubscriptionItems {
			if quantity, ok := newMap[item.Price.ID]; ok {
				//替换
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Price:    stripe.String(item.Price.ID),
					Quantity: stripe.Int64(quantity),
				})
				delete(newMap, item.Price.ID)
			} else {
				//删除之前全部，新增 Plan 和 Addons 方式
				items = append(items, &stripe.SubscriptionItemsParams{
					ID:       stripe.String(item.ID),
					Quantity: stripe.Int64(0),
				})
			}
		}
		//新增剩余的
		for channelPlanId, quantity := range newMap {
			items = append(items, &stripe.SubscriptionItemsParams{
				Price:    stripe.String(channelPlanId),
				Quantity: stripe.Int64(quantity),
			})
		}
	}

	return items, nil
}

// DoRemoteChannelSubscriptionUpdate 需保证同一个 Price 在 Items 中不能出现两份
func (s Stripe) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "out channel not found")
	stripe.Key = channelEntity.ChannelSecret
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
			params.PaymentBehavior = stripe.String("pending_if_incomplete") //pendingIfIncomplete 只有部分字段可以更新 Price Quantity
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice))
		}
	} else {
		if consts.NonEffectImmediatelyUsePendingUpdate {
			params.ProrationDate = stripe.Int64(subscriptionRo.ProrationDate)
			params.PaymentBehavior = stripe.String("pending_if_incomplete") //pendingIfIncomplete 只有部分字段可以更新 Price Quantity
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorAlwaysInvoice))
		} else {
			params.ProrationBehavior = stripe.String(string(stripe.SubscriptionSchedulePhaseProrationBehaviorNone))
		}
	}
	updateSubscription, err := sub.Update(subscriptionRo.Subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", params, updateSubscription, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}

	if subscriptionRo.EffectImmediate && !consts.ProrationUsingUniBeeCompute {
		queryParams := &stripe.InvoiceParams{}
		newInvoice, err := invoice.Get(updateSubscription.LatestInvoice.ID, queryParams)
		log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", queryParams, newInvoice, err, "GetInvoice", nil, channelEntity)
		g.Log().Infof(ctx, "query new invoice:", newInvoice)

		return &ro.ChannelUpdateSubscriptionInternalResp{
			Data:            utility.FormatToJsonString(updateSubscription),
			ChannelUpdateId: newInvoice.ID,
			Link:            newInvoice.HostedInvoiceURL,
			Paid:            newInvoice.Paid,
		}, nil
	} else {
		//EffectImmediate=false 不需要支付 获取的发票是之前最新的发票
		return &ro.ChannelUpdateSubscriptionInternalResp{
			Data: utility.FormatToJsonString(updateSubscription),
			Paid: false,
		}, nil
	}
}

// DoRemoteChannelSubscriptionDetails 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
func (s Stripe) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.SubscriptionParams{}
	response, err := sub.Get(subscription.ChannelSubscriptionId, params)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionDetails", params, response, err, subscription.ChannelSubscriptionId, nil, channelEntity)
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

// DoRemoteChannelPlanActive 使用 price 代替 plan  https://stripe.com/docs/api/plans
func (s Stripe) DoRemoteChannelPlanActive(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(true) // todo mark 使用这种方式可能不能用
	result, err := price.Update(planChannel.ChannelPlanId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanActive", params, result, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelPlanDeactivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(false) // todo mark 使用这种方式可能不能用
	result, err := price.Update(planChannel.ChannelPlanId, params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanDeactivate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.ProductParams{
		Active:      stripe.Bool(true),
		Description: stripe.String(plan.ChannelProductDescription), // todo mark 暂时不确定 description 如果为空会怎么样
		Name:        stripe.String(plan.ChannelProductName),
	}
	if len(plan.ImageUrl) > 0 {
		params.Images = stripe.StringSlice([]string{plan.ImageUrl})
	}
	if len(plan.HomeUrl) > 0 {
		params.URL = stripe.String(plan.HomeUrl)
	}
	result, err := product.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelProductCreate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	//Prod 创建好了之后似乎并不是Active 状态 todo mark
	return &ro.ChannelCreateProductInternalResp{
		ChannelProductId:     result.ID,
		ChannelProductStatus: fmt.Sprintf("%v", result.Active),
	}, nil
}

func (s Stripe) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.ChannelPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	stripe.Key = channelEntity.ChannelSecret
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
	// 使用 Price 代替 Plan https://stripe.com/docs/api/plans
	params := &stripe.PriceParams{
		Currency:   stripe.String(strings.ToLower(targetPlan.Currency)),
		UnitAmount: stripe.Int64(targetPlan.Amount),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(targetPlan.IntervalUnit),
			IntervalCount: stripe.Int64(int64(targetPlan.IntervalCount)),
		},
		Product: stripe.String(planChannel.ChannelProductId),
		Metadata: map[string]string{
			"PlanId": strconv.FormatUint(targetPlan.Id, 10),
			"Type":   strconv.Itoa(targetPlan.Type),
		},
		//ProductData: &stripe.PriceProductDataParams{
		//	ID:   stripe.String(planChannel.ChannelProductId),
		//	Name: stripe.String(targetPlan.PlanName),
		//},//这里是创建的意思
	}
	result, err := price.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelPlanCreateAndActivate", params, result, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCreatePlanInternalResp{
		ChannelPlanId:     result.ID,
		ChannelPlanStatus: fmt.Sprintf("%v", result.Active),
		Data:              utility.FormatToJsonString(result),
		Status:            consts.PlanChannelStatusActive,
	}, nil
}

func (s Stripe) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantChannelConfig, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoiceMarkUncollectibleParams{}
	response, err := invoice.MarkUncollectible(cancelInvoiceInternalReq.ChannelInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCancel", params, response, err, "", nil, payChannel)
	return parseStripeInvoice(response, int64(payChannel.Id)), nil
}

func (s Stripe) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantChannelConfig, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	channelUser := queryAndCreateChannelUser(ctx, payChannel, createInvoiceInternalReq.Invoice.UserId)

	params := &stripe.InvoiceParams{
		Currency: stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
		Customer: stripe.String(channelUser.ChannelUserId)}
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
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCreateAndPay", params, result, err, "New", nil, payChannel)

	for _, line := range createInvoiceInternalReq.InvoiceLines {
		ItemParams := &stripe.InvoiceItemParams{
			Invoice:  stripe.String(result.ID),
			Currency: stripe.String(strings.ToLower(createInvoiceInternalReq.Invoice.Currency)), //小写
			//UnitAmount:  stripe.Int64(line.UnitAmountExcludingTax),
			Amount:      stripe.Int64(line.Amount),
			Description: stripe.String(line.Description),
			//Quantity:    stripe.Int64(line.Quantity),
			Customer: stripe.String(channelUser.ChannelUserId)}
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
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceCreateAndPay", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, payChannel)
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

	return &ro.ChannelDetailInvoiceInternalResp{
		TotalAmount:                    detail.Total,
		TotalAmountExcludingTax:        detail.TotalExcludingTax,
		TaxAmount:                      detail.Tax,
		SubscriptionAmount:             detail.Subtotal,
		SubscriptionAmountExcludingTax: detail.TotalExcludingTax,
		Currency:                       strings.ToUpper(string(detail.Currency)),
		Lines:                          invoiceItems,
		ChannelId:                      int64(payChannel.Id),
		Status:                         status,
		ChannelUserId:                  detail.Customer.ID,
		Link:                           detail.HostedInvoiceURL,
		ChannelStatus:                  string(detail.Status),
		ChannelInvoiceId:               detail.ID,
		ChannelInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
	}, nil
}

func (s Stripe) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantChannelConfig, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.InvoicePayParams{}
	response, err := invoice.Pay(payInvoiceInternalReq.ChannelInvoiceId, params)
	log.SaveChannelHttpLog("DoRemoteChannelInvoicePay", params, response, err, "", nil, payChannel)
	return parseStripeInvoice(response, int64(payChannel.Id)), nil
}

func (s Stripe) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantChannelConfig, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	utility.Assert(payChannel != nil, "channel not found")
	stripe.Key = payChannel.ChannelSecret
	s.setUnibeeAppInfo()

	params := &stripe.InvoiceParams{}
	detail, err := invoice.Get(channelInvoiceId, params)
	if err != nil {
		return nil, err
	}
	log.SaveChannelHttpLog("DoRemoteChannelInvoiceDetails", params, detail, err, "", nil, payChannel)
	return parseStripeInvoice(detail, int64(payChannel.Id)), nil
}

func (s Stripe) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	utility.Assert(createPayContext.PayChannel != nil, "channel not found")
	stripe.Key = createPayContext.PayChannel.ChannelSecret
	s.setUnibeeAppInfo()
	channelUser := queryAndCreateChannelUser(ctx, createPayContext.PayChannel, createPayContext.Pay.UserId)

	if createPayContext.CheckoutMode {

		//channelVatRate := query.GetSubscriptionVatRateChannel(ctx, subscriptionRo.VatCountryRate.Id, createPayContext.PayChannel.Id)
		//if channelVatRate == nil {
		//	params := &stripe.TaxRateParams{
		//		DisplayName: stripe.String("VAT"),
		//		//Description: stripe.String(createPayContext.Pay.CountryName),
		//		Percentage: stripe.Float64(utility.ConvertTaxScaleToPercentageFloat(createPayContext.Pay.StandardTaxPercentage)),
		//		Country:    stripe.String(createPayContext.Pay.CountryCode),
		//		Active:     stripe.Bool(true),
		//		//Jurisdiction: stripe.String("DE"),
		//		Inclusive: stripe.Bool(false),
		//	}
		//	vatCreateResult, err := taxrate.New(params)
		//	if err != nil {
		//		g.Log().Printf(ctx, "taxrate.New: %v", err.Error())
		//		return nil, err
		//	}
		//	channelVatRate = &entity.ChannelVatRate{
		//		VatRateId:        int64(subscriptionRo.VatCountryRate.Id),
		//		ChannelId:        int64(createPayContext.PayChannel.Id),
		//		ChannelVatRateId: vatCreateResult.ID,
		//	}
		//	result, err := dao.ChannelVatRate.Ctx(ctx).Data(channelVatRate).OmitNil().Insert(channelVatRate)
		//	if err != nil {
		//		err = gerror.Newf(`SubscriptionVatRateChannel record insert failure %s`, err.Error())
		//		return nil, err
		//	}
		//	id, _ := result.LastInsertId()
		//	channelVatRate.Id = uint64(uint(id))
		//}

		var items []*stripe.CheckoutSessionLineItemParams
		for _, line := range createPayContext.Invoice.Lines {
			items = append(items, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(strings.ToLower(line.Currency)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(line.Description),
					},
					UnitAmount: stripe.Int64(line.UnitAmountExcludingTax),
				},
				Quantity: stripe.Int64(line.Quantity),
				//TaxRates: []*string{stripe.String(channelVatRate.ChannelVatRateId)}, // todo mark tax add
			})
		}
		checkoutParams := &stripe.CheckoutSessionParams{
			Customer:  stripe.String(channelUser.ChannelUserId),
			Currency:  stripe.String(strings.ToLower(createPayContext.Pay.Currency)),
			LineItems: items,
			//AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			//	Enabled: stripe.Bool(false), //Default值 false，表示不需要 stripe 计算税率，true 反之
			//},
			SuccessURL: stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, true)),
			CancelURL:  stripe.String(webhook2.GetPaymentRedirectEntranceUrlCheckout(createPayContext.Pay, false)),
		}
		checkoutParams.Mode = stripe.String(string(stripe.CheckoutSessionModePayment))
		checkoutParams.Metadata = createPayContext.MediaData
		//checkoutParams.ExpiresAt
		detail, err := session.New(checkoutParams)
		if err != nil {
			return nil, err
		}
		log.SaveChannelHttpLog("DoRemoteChannelPayment", checkoutParams, detail, err, "CheckoutSession", nil, createPayContext.PayChannel)
		var status consts.PayStatusEnum = consts.TO_BE_PAID
		if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusOpen)) == 0 {
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusComplete)) == 0 {
			status = consts.PAY_SUCCESS
		} else if strings.Compare(string(detail.Status), string(stripe.CheckoutSessionStatusExpired)) == 0 {
			status = consts.PAY_FAILED
		}
		var channelPaymentId string
		if detail.PaymentIntent != nil {
			channelPaymentId = detail.PaymentIntent.ID
		}
		return &ro.CreatePayInternalResp{
			Status:                 status,
			ChannelPaymentId:       channelPaymentId,
			ChannelPaymentIntentId: detail.ID,
			Link:                   detail.URL,
		}, nil
	} else {
		params := &stripe.InvoiceParams{
			Metadata: createPayContext.MediaData,
			Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
			Customer: stripe.String(channelUser.ChannelUserId)}

		if createPayContext.PayMethod == 1 {
			params.CollectionMethod = stripe.String("charge_automatically")
			if len(createPayContext.ChannelPaymentMethod) > 0 {
				params.DefaultPaymentMethod = stripe.String(createPayContext.ChannelPaymentMethod)
			} else if len(channelUser.ChannelDefaultPaymentMethod) > 0 {
				params.DefaultPaymentMethod = stripe.String(channelUser.ChannelDefaultPaymentMethod)
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
		log.SaveChannelHttpLog("DoRemoteChannelPayment", params, result, err, "New", nil, createPayContext.PayChannel)

		for _, line := range createPayContext.Invoice.Lines {
			ItemParams := &stripe.InvoiceItemParams{
				Invoice:  stripe.String(result.ID),
				Currency: stripe.String(strings.ToLower(createPayContext.Invoice.Currency)),
				//UnitAmount:  stripe.Int64(line.UnitAmountExcludingTax),
				Amount:      stripe.Int64(line.Amount),
				Description: stripe.String(line.Description),
				//Quantity:    stripe.Int64(line.Quantity),
				Customer: stripe.String(channelUser.ChannelUserId)}
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
		log.SaveChannelHttpLog("DoRemoteChannelPayment", finalizeInvoiceParam, detail, err, "FinalizeInvoice", nil, createPayContext.PayChannel)
		if err != nil {
			return nil, err
		}
		if createPayContext.PayImmediate {
			params := &stripe.InvoicePayParams{}
			response, err := invoice.Pay(result.ID, params)
			log.SaveChannelHttpLog("DoRemoteChannelPayment", params, response, err, "PayInvoice", nil, createPayContext.PayChannel)
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
		var channelPaymentId string
		if detail.PaymentIntent != nil {
			channelPaymentId = detail.PaymentIntent.ID
		}
		return &ro.CreatePayInternalResp{
			Status:                 status,
			ChannelPaymentId:       channelPaymentId,
			ChannelPaymentIntentId: detail.ID,
			Link:                   detail.HostedInvoiceURL,
		}, nil
	}
}

func (s Stripe) DoRemoteChannelCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	utility.Assert(payment.ChannelId > 0, "invalid payment channelId")
	utility.Assert(len(payment.ChannelPaymentIntentId) > 0, "invalid payment ChannelPaymentIntentId")
	channelEntity := util.GetOverseaPayChannel(ctx, payment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	s.setUnibeeAppInfo()
	params := &stripe.InvoiceVoidInvoiceParams{}
	result, err := invoice.VoidInvoice(payment.ChannelPaymentIntentId, params)
	if err != nil {
		return nil, err
	}
	invoiceDetails := parseStripeInvoice(result, payment.ChannelId)
	var status = consts.TO_BE_PAID
	if invoiceDetails.Status == consts.InvoiceStatusPaid {
		status = consts.PAY_SUCCESS
	} else if invoiceDetails.Status == consts.InvoiceStatusFailed || invoiceDetails.Status == consts.InvoiceStatusCancelled {
		status = consts.PAY_FAILED
	}
	return &ro.OutPayCancelRo{
		MerchantId:      strconv.FormatInt(payment.MerchantId, 10),
		ChannelCancelId: result.ID,
		Reference:       payment.PaymentId,
		Status:          strconv.Itoa(status),
	}, nil
}

func (s Stripe) DoRemoteChannelPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefund(ctx context.Context, payment *entity.Payment, one *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	utility.Assert(payment.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, payment.ChannelId)
	utility.Assert(channelEntity != nil, "channel not found")
	params := &stripe.RefundParams{PaymentIntent: stripe.String(payment.ChannelPaymentId)}
	params.Reason = stripe.String(one.RefundComment)
	params.Amount = stripe.Int64(one.RefundAmount)
	params.Currency = stripe.String(strings.ToLower(one.Currency))
	params.Metadata = map[string]string{"RefundId": one.RefundId}
	result, err := refund.New(params)
	log.SaveChannelHttpLog("DoRemoteChannelRefund", params, result, err, "refund", nil, channelEntity)
	utility.Assert(err == nil, fmt.Sprintf("call stripe refund error %s", err))
	utility.Assert(result != nil, "Stripe refund failed, result is nil")
	return &ro.OutPayRefundRo{
		ChannelRefundId: result.ID,
		Status:          consts.REFUND_ING,
	}, nil
}

func parseStripeRefund(item *stripe.Refund) *ro.OutPayRefundRo {
	var channelPaymentId string
	if item.PaymentIntent != nil {
		channelPaymentId = item.PaymentIntent.ID
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
		ChannelRefundId:  item.ID,
		ChannelPaymentId: channelPaymentId,
		Status:           status,
		Reason:           string(item.Reason),
		RefundFee:        item.Amount,
		Currency:         strings.ToUpper(string(item.Currency)),
		RefundTime:       gtime.NewFromTimeStamp(item.Created),
	}
}

func parseStripePayment(item *stripe.PaymentIntent, payChannel *entity.MerchantChannelConfig) *ro.ChannelPaymentRo {
	var channelInvoiceId string
	if item.Invoice != nil {
		channelInvoiceId = item.Invoice.ID
	}
	var channelUserId string
	if item.Customer != nil {
		channelUserId = item.Customer.ID
	}
	var status = consts.TO_BE_PAID
	if strings.Compare(string(item.Status), "succeeded") == 0 {
		status = consts.PAY_SUCCESS
	} else if strings.Compare(string(item.Status), "canceled") == 0 {
		status = consts.PAY_CANCEL
	}
	var captureStatus = consts.WAITING_AUTHORIZED
	if strings.Compare(string(item.Status), "requires_capture") == 0 {
		captureStatus = consts.AUTHORIZED
	} else if strings.Compare(string(item.Status), "requires_confirmation") == 0 {
		captureStatus = consts.CAPTURE_REQUEST
	}
	return &ro.ChannelPaymentRo{
		ChannelId:        int64(payChannel.Id),
		MerchantId:       payChannel.MerchantId,
		ChannelInvoiceId: channelInvoiceId,
		ChannelUserId:    channelUserId,
		ChannelPaymentId: item.ID,
		Status:           status,
		CaptureStatus:    captureStatus,
		TotalAmount:      item.Amount,
		PaymentAmount:    item.AmountReceived,
		Currency:         strings.ToUpper(string(item.Currency)),
		PayTime:          gtime.NewFromTimeStamp(item.Created),
		CreateTime:       gtime.NewFromTimeStamp(item.Created),
		CancelTime:       gtime.NewFromTimeStamp(item.CanceledAt),
		CancelReason:     string(item.CancellationReason),
	}
}

func parseStripeSubscription(subscription *stripe.Subscription) *ro.ChannelDetailSubscriptionInternalResp {
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
	var channelDefaultPaymentMethod = ""
	if subscription.DefaultPaymentMethod != nil {
		channelDefaultPaymentMethod = subscription.DefaultPaymentMethod.ID
	}

	return &ro.ChannelDetailSubscriptionInternalResp{
		Status:                      status,
		ChannelSubscriptionId:       subscription.ID,
		ChannelStatus:               string(subscription.Status),
		Data:                        utility.FormatToJsonString(subscription),
		ChannelItemData:             utility.MarshalToJsonString(subscription.Items.Data),
		ChannelLatestInvoiceId:      subscription.LatestInvoice.ID,
		ChannelLatestPaymentId:      latestChannelPaymentId,
		ChannelDefaultPaymentMethod: channelDefaultPaymentMethod,
		CancelAtPeriodEnd:           subscription.CancelAtPeriodEnd,
		CurrentPeriodStart:          subscription.CurrentPeriodStart,
		CurrentPeriodEnd:            subscription.CurrentPeriodEnd,
		BillingCycleAnchor:          subscription.BillingCycleAnchor,
		TrialEnd:                    subscription.TrialEnd,
	}
}

func parseStripeInvoice(detail *stripe.Invoice, channelId int64) *ro.ChannelDetailInvoiceInternalResp {
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

	var channelPaymentId string
	if detail.PaymentIntent != nil {
		channelPaymentId = detail.PaymentIntent.ID
	}
	var channelSubscriptionId string
	if detail.Subscription != nil {
		channelSubscriptionId = detail.Subscription.ID
	}
	var subscriptionId string
	if detail.SubscriptionDetails != nil {
		subscriptionId = detail.SubscriptionDetails.Metadata["SubId"]
	}
	var channelUserId string
	if detail.Customer != nil {
		channelUserId = detail.Customer.ID
	}
	var paymentTime int64
	var cancelTime int64
	if detail.StatusTransitions != nil {
		paymentTime = detail.StatusTransitions.PaidAt
		cancelTime = detail.StatusTransitions.VoidedAt
	}
	var channelDefaultPaymentMethod = ""
	if detail.DefaultPaymentMethod != nil {
		channelDefaultPaymentMethod = detail.DefaultPaymentMethod.ID
	}

	return &ro.ChannelDetailInvoiceInternalResp{
		ChannelDefaultPaymentMethod:    channelDefaultPaymentMethod,
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
		ChannelId:                      channelId,
		Status:                         status,
		Link:                           detail.HostedInvoiceURL,
		ChannelStatus:                  string(detail.Status),
		ChannelInvoicePdf:              detail.InvoicePDF,
		PeriodStart:                    detail.PeriodStart,
		PeriodEnd:                      detail.PeriodEnd,
		ChannelInvoiceId:               detail.ID,
		ChannelUserId:                  channelUserId,
		ChannelSubscriptionId:          channelSubscriptionId,
		SubscriptionId:                 subscriptionId,
		ChannelPaymentId:               channelPaymentId,
		PaymentTime:                    paymentTime,
		Reason:                         string(detail.BillingReason),
		CreateTime:                     detail.Created,
		CancelTime:                     cancelTime,
	}
}
