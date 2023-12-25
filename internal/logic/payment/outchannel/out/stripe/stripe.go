package stripe

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"log"
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
	if len(subscriptionRo.Subscription.ChannelUserId) == 0 {
		params := &stripe.CustomerParams{
			// todo mark 创建 customer 需要这两个字段 https://stripe.com/docs/api/customers/create
			Name:  stripe.String(subscriptionRo.Subscription.CustomerName),
			Email: stripe.String(subscriptionRo.Subscription.CustomerEmail),
		}

		createCustomResult, err := customer.New(params)
		if err != nil {
			log.Printf("customer.New: %v", err)
			return nil, err
		}
		subscriptionRo.Subscription.ChannelUserId = createCustomResult.ID
	}
	taxInclusive := true
	if subscriptionRo.Plan.TaxInclusive == 0 {
		//税费不包含
		taxInclusive = false
	}
	//stripe.CheckoutSessionParams{
	//	Params:                   stripe.Params{},
	//	AfterExpiration:          nil,
	//	AllowPromotionCodes:      nil,
	//	AutomaticTax:             nil,
	//	BillingAddressCollection: nil,
	//	CancelURL:                nil,
	//	ClientReferenceID:        nil,
	//	ConsentCollection:        nil,
	//	Currency:                 nil,
	//	Customer:                 nil,
	//	CustomerCreation:         nil,
	//	CustomerEmail:            nil,
	//	CustomerUpdate:           nil,
	//	Discounts:                nil,
	//	ExpiresAt:                nil,
	//	LineItems:                nil,
	//	Locale:                   nil,
	//	Mode:                     nil,
	//	PaymentIntentData:        nil,
	//	PaymentMethodOptions:     nil,
	//	PaymentMethodTypes:       nil,
	//	PhoneNumberCollection:    nil,
	//	SetupIntentData:          nil,
	//	SubmitType:               nil,
	//	SubscriptionData:         nil,
	//	SuccessURL:               nil,
	//	TaxIDCollection:          nil,
	//}
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(subscriptionRo.Subscription.ChannelUserId),
		Currency: stripe.String(strings.ToLower(subscriptionRo.Plan.Currency)), //小写
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(subscriptionRo.PlanChannel.ChannelPlanId),
				//TaxRates: stripe.TaxRate{ // todo mark 取消 stripe 计算费率
				//	APIResource:  stripe.APIResource{},
				//	Active:       false,
				//	Country:      "",
				//	Created:      0,
				//	Description:  "",
				//	DisplayName:  "",
				//	ID:           "",
				//	Inclusive:    false,
				//	Jurisdiction: "",
				//	Livemode:     false,
				//	Metadata:     nil,
				//	Object:       "",
				//	Percentage:   0,
				//	State:        "",
				//	TaxType:      "",
				//},
			},
		},
		AutomaticTax: &stripe.SubscriptionAutomaticTaxParams{
			Enabled: stripe.Bool(!taxInclusive), //默认值 false，表示不需要 stripe 计算税率，true 反之 todo 添加 item 里面的 tax_tates
		},
		PaymentBehavior:  stripe.String("default_incomplete"),   // todo mark https://stripe.com/docs/api/subscriptions/create
		CollectionMethod: stripe.String("charge_automatically"), //默认行为 charge_automatically，自动扣款
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	createSubscription, err := sub.New(subscriptionParams)
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(createSubscription)
	return &ro.CreateSubscriptionInternalResp{
		ChannelSubscriptionId:     createSubscription.ID,
		ChannelSubscriptionStatus: string(createSubscription.Status),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil
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

func (s Stripe) DoRemoteChannelSubscriptionWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
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

func (s Stripe) DoRemoteChannelWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelRedirect(r *ghttp.Request) {
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
