package stripe

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/plan"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"log"
	"strconv"
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

func (s Stripe) DoRemoteChannelSubscriptionCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CreateSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, uint64(planChannel.ChannelId))
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	if len(subscription.ChannelUserId) == 0 {
		params := &stripe.CustomerParams{
			// todo mark 创建 customer 需要这两个字段 https://stripe.com/docs/api/customers/create
			Name:  stripe.String(strconv.FormatInt(subscription.UserId, 10)),
			Email: stripe.String(strconv.FormatInt(subscription.UserId, 10)),
		}

		createCustomResult, err := customer.New(params)
		if err != nil {
			log.Printf("customer.New: %v", err)
			return nil, err
		}
		subscription.ChannelUserId = createCustomResult.ID
	}
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(subscription.ChannelUserId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				// todo mark
				//Price: stripe.String(req.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
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
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelSubscriptionUpdate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.UpdateSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ListSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelSubscriptionWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (s Stripe) DoRemoteChannelPlanActive(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, uint64(planChannel.ChannelId))
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PlanParams{}
	params.Active = stripe.Bool(true) // todo mark 使用这种方式可能不能用
	_, err = plan.Update(planChannel.ChannelPlanId, params)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelPlanDeactivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, uint64(planChannel.ChannelId))
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PlanParams{}
	params.Active = stripe.Bool(false) // todo mark 使用这种方式可能不能用
	_, err = plan.Update(planChannel.ChannelPlanId, params)
	if err != nil {
		return err
	}
	return nil
}

func (s Stripe) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, uint64(planChannel.ChannelId))
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.ProductParams{
		Active:      stripe.Bool(true),
		Description: stripe.String(plan.ChannelProductDescription),
		Images:      stripe.StringSlice([]string{plan.ImageUrl}),
		Name:        stripe.String(plan.ChannelProductName),
		URL:         stripe.String(plan.HomeUrl),
	}
	result, err := product.New(params)
	if err != nil {
		return nil, err
	}
	return &ro.CreateProductInternalResp{
		ChannelProductId:     result.ID,
		ChannelProductStatus: fmt.Sprintf("%v", result.Active),
	}, nil
}

func (s Stripe) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, targetPlan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, uint64(planChannel.ChannelId))
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	stripe.Key = channelEntity.ChannelSecret
	s.setUnibeeAppInfo()
	params := &stripe.PlanParams{
		//todo mark
		Active:   stripe.Bool(true),
		Amount:   stripe.Int64(1200),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Interval: stripe.String(string(stripe.PlanIntervalMonth)),
		Product:  &stripe.PlanProductParams{ID: stripe.String("prod_NjpI7DbZx6AlWQ")},
	}
	result, err := plan.New(params)
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
