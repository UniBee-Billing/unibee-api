package paypal

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/plutov/paypal/v4"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
	"net/http"
	"time"
)

// 接口文档：https://developer.paypal.com/docs/api/payments/v1/#payment_create
// https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
// clientId ATaWQ8G9oJNFyle9YCt59
// Secret EHUy5GALkYr1Qp0n6MepJY8LnUwYCBIWElG4Iv_DO3mdYcbB2l6zwJxk99OrPhbdNRLk7GkHEqb5RHEA

// Other ClientId AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur
// Other Secret EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw

//APIBaseSandBox = "https://api-m.sandbox.paypal.com"
//APIBaseLive = "https://api-m.paypal.com"

type Paypal struct {
}

func init() {
	//注册 webhooks
}

// todo mark 确认改造成单例是否可行，不用每次都去获取 accessToken
func NewClient(clientID string, secret string, APIBase string) (*paypal.Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &paypal.Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

func (p Paypal) DoRemoteChannelSubscriptionCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CreateSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	createSubscription, err := client.CreateSubscription(ctx, paypal.SubscriptionBase{
		PlanID: planChannel.ChannelPlanId,
		// todo mark
		StartTime:          nil,
		EffectiveTime:      nil,
		Quantity:           "",
		ShippingAmount:     nil,
		Subscriber:         nil,
		AutoRenewal:        false,
		ApplicationContext: nil,
		CustomID:           "",
		Plan:               nil,
	})
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(createSubscription)
	return &ro.CreateSubscriptionInternalResp{
		ChannelSubscriptionId:     createSubscription.ID,
		ChannelSubscriptionStatus: string(createSubscription.SubscriptionStatus),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil
}

// todo mark paypal 的 cancel 似乎是无法恢复的，和 stripe 不一样，需要确认是否有真实 cancel 的需求
func (p Paypal) DoRemoteChannelSubscriptionCancel(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.CancelSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	err = client.CancelSubscription(ctx, subscription.ChannelSubscriptionId, "")
	if err != nil {
		return nil, err
	} // cancelReason

	return &ro.CancelSubscriptionInternalResp{}, nil //todo mark
}

func (p Paypal) DoRemoteChannelSubscriptionUpdate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.UpdateSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	err = client.UpdateSubscription(ctx, paypal.Subscription{}) //todo mark
	if err != nil {
		return nil, err
	}

	return &ro.UpdateSubscriptionInternalResp{}, nil //todo mark
}

func (p Paypal) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ListSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	_, err = client.GetSubscriptionDetails(ctx, subscription.ChannelSubscriptionId)
	if err != nil {
		return nil, err
	}

	return nil, nil //todo mark
}

func (p Paypal) DoRemoteChannelSubscriptionWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.ActivateSubscriptionPlan(ctx, planChannel.ChannelPlanId)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.DeactivateSubscriptionPlans(ctx, planChannel.ChannelPlanId)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 outchannel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	productResult, err := client.CreateProduct(ctx, paypal.Product{
		Name:        plan.ChannelProductName,
		Description: plan.ChannelProductDescription,
		Category:    paypal.ProductCategorySoftware,
		Type:        paypal.ProductTypeService,
		ImageUrl:    plan.ImageUrl,
		HomeUrl:     plan.HomeUrl,
	})
	if err != nil {
		return nil, err
	}
	return &ro.CreateProductInternalResp{
		ChannelProductId:     productResult.ID,
		ChannelProductStatus: "",
	}, nil
}

func (p Paypal) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.CreatePlanInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	subscriptionPlan, err := client.CreateSubscriptionPlan(ctx, paypal.SubscriptionPlan{
		ProductId:   planChannel.ChannelProductId,
		Name:        plan.PlanName,
		Status:      paypal.SubscriptionPlanStatusActive,
		Description: plan.Description,
		//todo mark
		BillingCycles: []paypal.BillingCycle{
			{
				PricingScheme: paypal.PricingScheme{
					Version: 1,
					FixedPrice: paypal.Money{
						Currency: "EUR",
						Value:    "5",
					},
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				},
				Frequency: paypal.Frequency{
					IntervalUnit:  paypal.IntervalUnitYear,
					IntervalCount: 1,
				},
				TenureType:  paypal.TenureTypeRegular,
				Sequence:    1,
				TotalCycles: 0,
			},
		},
		PaymentPreferences: &paypal.PaymentPreferences{
			AutoBillOutstanding:     false,
			SetupFee:                nil,
			SetupFeeFailureAction:   paypal.SetupFeeFailureActionCancel,
			PaymentFailureThreshold: 0,
		},
		Taxes: &paypal.Taxes{
			Percentage: "19",
			Inclusive:  false,
		},
		QuantitySupported: false,
	})
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(subscriptionPlan)
	return &ro.CreatePlanInternalResp{
		ChannelPlanId:     subscriptionPlan.ID,
		ChannelPlanStatus: string(subscriptionPlan.Status),
		Data:              string(jsonData),
		Status:            consts.PlanStatusActive,
	}, nil
}

func (p Paypal) DoRemoteChannelWebhook(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRedirect(r *ghttp.Request) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCapture(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCancel(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPayStatusCheck(ctx context.Context, pay *entity.OverseaPay) (res *ro.OutPayRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefundStatusCheck(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefund(ctx context.Context, pay *entity.OverseaPay, refund *entity.OverseaRefund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
