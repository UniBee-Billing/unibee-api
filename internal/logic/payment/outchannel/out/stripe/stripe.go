package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	sub "github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/out"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"net/http"
	"strings"
)

type Stripe struct {
}

func (s Stripe) setUnibeeAppInfo() {
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "unibee.server",
		Version: "0.0.1",
		URL:     "https://unibee.dev",
	})
}

func (s Stripe) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.CreateSubscriptionRo) (res *ro.CreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	//{
	//	if len(subscriptionRo.Subscription.ChannelUserId) == 0 {
	//		params := &stripe.CustomerParams{
	//			// todo mark 创建 customer 需要这两个字段 https://stripe.com/docs/api/customers/create
	//			Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
	//			Email: stripe.String(subscriptionRo.Subscription.CustomerEmail),
	//		}
	//
	//		createCustomResult, err := customer.New(params)
	//		if err != nil {
	//			log.Printf("customer.New: %v", err)
	//			return nil, err
	//		}
	//		subscriptionRo.Subscription.ChannelUserId = createCustomResult.ID
	//	}
	//	taxInclusive := true
	//	if subscriptionRo.Plan.TaxInclusive == 0 {
	//		//税费不包含
	//		taxInclusive = false
	//	}
	//	subscriptionParams := &stripe.SubscriptionParams{
	//		Customer: stripe.String(subscriptionRo.Subscription.ChannelUserId),
	//		Currency: stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
	//		Items: []*stripe.SubscriptionItemsParams{
	//			{
	//				Price: stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
	//			},
	//		},
	//		AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
	//			Enabled: stripe.Bool(!taxInclusive), //默认值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_tates
	//		},
	//		PaymentBehavior:  stripe.String("default_incomplete"),   // todo mark https://stripe.com/docs/api/subscriptions/create
	//		CollectionMethod: stripe.String("charge_automatically"), //默认行为 charge_automatically，自动扣款
	//	}
	//	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	//	createSubscription, err := sub.New(subscriptionParams)
	//	if err != nil {
	//		return nil, err
	//	}
	//	////尝试创建发票
	//	//params := &stripe.InvoiceParams{
	//	//	Customer:     stripe.String(subscriptionRo.Subscription.ChannelUserId),
	//	//	Subscription: stripe.String(createSubscription.ID),
	//	//}
	//	//createInvoice, err := invoice.New(params)
	//	//if err != nil {
	//	//	return nil, err
	//	//}
	//	//createPayInvoice, err := invoice.Pay(createInvoice.ID, &stripe.InvoicePayParams{})
	//	//if err != nil {
	//	//	return nil, err
	//	//}
	//	//createPayInvoiceJson, _ := gjson.Marshal(createPayInvoice)
	//	//g.Log().Infof(ctx, "pay invoice:%s", string(createPayInvoiceJson))
	//
	//	jsonData, _ := gjson.Marshal(createSubscription)
	//	return &ro.CreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     createSubscription.ID,
	//		ChannelSubscriptionStatus: string(createSubscription.Status),
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}
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
	//	return &ro.CreateSubscriptionInternalResp{
	//		ChannelSubscriptionId:     createSubscription.ID,
	//		ChannelSubscriptionStatus: "true",
	//		Data:                      string(jsonData),
	//		Status:                    0, //todo mark
	//	}, nil
	//}
	{
		checkoutParams := &stripe.CheckoutSessionParams{
			Metadata: map[string]string{
				"orderId": subscriptionRo.Subscription.SubscriptionId,
			},
			Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
					Quantity: stripe.Int64(1),
				},
			},
			SuccessURL: stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, true)),
			CancelURL:  stripe.String(out.GetSubscriptionRedirectEntranceUrl(subscriptionRo.Subscription, false)),
		}

		result, err := session.New(checkoutParams)
		if err != nil {
			return nil, err
		}
		jsonData, _ := gjson.Marshal(result)
		return &ro.CreateSubscriptionInternalResp{
			ChannelSubscriptionId:     result.ID,
			ChannelSubscriptionStatus: "true",
			Data:                      string(jsonData),
			Status:                    0, //todo mark
		}, nil
	}

}

func (s Stripe) DoRemoteChannelSubscriptionCancel(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CancelSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.SubscriptionCancelParams{}
	_, err = sub.Cancel(subscription.ChannelSubscriptionId, params)
	if err != nil {
		return nil, err
	}
	return &ro.CancelSubscriptionInternalResp{}, nil //todo mark
}

// DoRemoteChannelSubscriptionUpdate 需保证同一个 Price 在 Items 中不能出现两份
func (s Stripe) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.UpdateSubscriptionRo) (res *ro.UpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()

	detail, err := sub.Get(subscriptionRo.Subscription.ChannelSubscriptionId, &stripe.SubscriptionParams{})
	if err != nil {
		return nil, err
	}
	//遍历
	var targetItems []*stripe.SubscriptionItemsParams
	for _, item := range detail.Items.Data {
		if strings.Compare(item.Price.ID, subscriptionRo.OldPlanChannel.ChannelPlanId) == 0 {
			targetItems = append(targetItems, &stripe.SubscriptionItemsParams{
				ID:    stripe.String(item.ID),
				Price: stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
			})
		}
	}
	if len(targetItems) == 0 {
		return nil, gerror.New("items not match")
	}

	params := &stripe.SubscriptionParams{
		Items: targetItems,
	}
	updateSubscription, err := sub.Update(subscriptionRo.Subscription.ChannelSubscriptionId, params)
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(updateSubscription)
	return &ro.UpdateSubscriptionInternalResp{
		ChannelSubscriptionId:     updateSubscription.ID,
		ChannelSubscriptionStatus: string(updateSubscription.Status),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil //todo mark
}

func (s Stripe) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ListSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.SubscriptionParams{}
	_, err = sub.Get("sub_1MowQVLkdIwHu7ixeRlqHVzs", params)
	if err != nil {
		return nil, err
	}
	return nil, nil //todo mark
}

// DoRemoteChannelCheckAndSetupWebhook https://stripe.com/docs/billing/subscriptions/webhooks
func (s Stripe) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	utility.Assert(payChannel != nil, "payChannel is nil")
	stripe.Key = payChannel.ChannelSecret
	params := &stripe.WebhookEndpointListParams{}
	params.Limit = stripe.Int64(10)
	result := webhookendpoint.List(params)
	if len(result.WebhookEndpointList().Data) > 1 {
		return gerror.New("webhook endpoints count > 1")
	}
	//过滤不可用
	if len(result.WebhookEndpointList().Data) == 0 {
		//创建
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				stripe.String("charge.succeeded"),
				stripe.String("charge.failed"),
			},
			URL: stripe.String(out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
		}
		result, err := webhookendpoint.New(params)
		if err != nil {
			return nil
		}
		//更新 secret
		utility.Assert(len(result.Secret) > 0, "secret is nil")
		err = query.UpdatePayChannelWebhookSecret(ctx, int64(payChannel.Id), result.Secret)
		if err != nil {
			return err
		}
	} else {
		utility.Assert(len(result.WebhookEndpointList().Data) == 1, "internal webhook update, count is not 1")
		//检查并更新, todo mark 优化检查逻辑，如果 evert 一致不用发起更新
		webhook := result.WebhookEndpointList().Data[0]
		utility.Assert(strings.Compare(webhook.Status, "enabled") == 0, "webhook not status enabled")
		params := &stripe.WebhookEndpointParams{
			EnabledEvents: []*string{
				//订阅相关 webhook
				stripe.String("customer.subscription.deleted"),
				stripe.String("customer.subscription.updated"),
				stripe.String("customer.subscription.created"),
				stripe.String("customer.subscription.trial_will_end"),
				stripe.String("customer.subscription.paused"),
				stripe.String("customer.subscription.resumed"),
			},
			URL: stripe.String(out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id))),
		}
		result, err := webhookendpoint.Update(webhook.ID, params)
		if err != nil {
			return err
		}
		utility.Assert(strings.Compare(result.Status, "enabled") == 0, "webhook not status enabled after updated")
	}

	return nil
}

// DoRemoteChannelPlanActive 使用 price 代替 plan  https://stripe.com/docs/api/plans
func (s Stripe) DoRemoteChannelPlanActive(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(true) // todo mark 使用这种方式可能不能用
	_, err = price.Update(planChannel.ChannelPlanId, params)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelPlanDeactivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PriceParams{}
	params.Active = stripe.Bool(false) // todo mark 使用这种方式可能不能用
	_, err = price.Update(planChannel.ChannelPlanId, params)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
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
	if err != nil {
		return nil, err
	}
	//Prod 创建好了之后似乎并不是Active 状态 todo mark
	return &ro.CreateProductInternalResp{
		ChannelProductId:     result.ID,
		ChannelProductStatus: fmt.Sprintf("%v", result.Active),
	}, nil
}

func (s Stripe) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
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
		UnitAmount: stripe.Int64(targetPlan.Amount), //todo mark 小数点可能不用处理
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String(targetPlan.IntervalUnit),
		},
		Product: stripe.String(planChannel.ChannelProductId),

		//ProductData: &stripe.PriceProductDataParams{
		//	ID:   stripe.String(planChannel.ChannelProductId),
		//	Name: stripe.String(targetPlan.PlanName),
		//},//这里是创建的意思
	}
	result, err := price.New(params)
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(result)
	return &ro.CreatePlanInternalResp{
		ChannelPlanId:     result.ID,
		ChannelPlanStatus: fmt.Sprintf("%v", result.Active),
		Data:              string(jsonData),
		Status:            consts.PlanStatusActive,
	}, nil
}

func (s Stripe) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	endpointSecret := payChannel.WebhookSecret
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(r.GetBody(), signatureHeader, endpointSecret)
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %v\n", payChannel.Channel, err)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	data, _ := gjson.Marshal(event)
	g.Log().Info(r.Context(), "Receive_Webhook_Channel:%s hook:%s", payChannel.Channel, string(data))

	switch event.Type {
	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err)
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}
		g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription deleted for %d.", payChannel.Channel, subscription.ID)
		// Then define and call a func to handle the deleted subscription.
		// handleSubscriptionCanceled(subscription)
	case "customer.subscription.updated":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err)
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}
		g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription updated for %s.", payChannel.Channel, subscription)
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handleSubscriptionUpdated(subscription)
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err)
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}
		g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription created for %s.", payChannel.Channel, subscription)
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handleSubscriptionCreated(subscription)
	case "customer.subscription.trial_will_end":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook JSON: %v\n", payChannel.Channel, err)
			r.Response.WriteHeader(http.StatusBadRequest)
			return
		}
		g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription trial will end for %d.", payChannel.Channel, subscription.ID)
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handleSubscriptionTrialWillEnd(subscription)
	default:
		g.Log().Errorf(r.Context(), "Webhook Channel:%s, Unhandled event type: %s\n", payChannel.Channel, event.Type)
	}
	r.Response.WriteHeader(http.StatusOK)
}

func (s Stripe) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
