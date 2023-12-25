package paypal

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/plutov/paypal/v4"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/out"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 接口文档：https://developer.paypal.com/docs/api/payments/v1/#payment_create
// https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
// clientId ATaWQ8G9oJNFyle9YCt59
// Secret EHUy5GALkYr1Qp0n6MepJY8LnUwYCBIWElG4Iv_DO3mdYcbB2l6zwJxk99OrPhbdNRLk7GkHEqb5RHEA

// Other ClientId AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur
// Other Secret EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw

// Other 2 ClientId AT-HU_WUeHCis_uqkU2Y8-0f54qq_QkoNXJeBj1-4S01__m1OLQn1jXnG9F86bcaH5TbcYiFed7UBRGH
// Other 2 Secret  EL2TLXWp_6XyZEtYqeRjLLVb9S_uYjwZOrBUiqhHhw96-50VisMsQvBDA09qMVntXrPf6TukiyfRCkG0

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

func (p Paypal) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.CreateSubscriptionRo) (res *ro.CreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	createSubscription, err := client.CreateSubscription(ctx, paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.ChannelPlanId,
		// todo mark
		StartTime:     nil,
		EffectiveTime: nil,
		Quantity:      "",
		//Plan: &paypal.PlanOverride{
		//	BillingCycles:      nil,
		//	PaymentPreferences: nil,
		//	Taxes:              nil,
		//},
		Subscriber:         nil,
		AutoRenewal:        false,
		ApplicationContext: nil,
		CustomID:           "",
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

// DoRemoteChannelSubscriptionUpdate 新旧 Plan 需要在同一个 Product 下，你这个 Product 有什么用，stripe 不需要
// 需要支付之后才能更新，stripe 不需要
func (p Paypal) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.UpdateSubscriptionRo) (res *ro.UpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	updateSubscription, err := client.ReviseSubscription(ctx, subscriptionRo.Subscription.ChannelSubscriptionId, paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.ChannelPlanId,
		//todo mark
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(updateSubscription)
	return &ro.UpdateSubscriptionInternalResp{
		ChannelSubscriptionId:     updateSubscription.ID,
		ChannelSubscriptionStatus: string(updateSubscription.SubscriptionStatus),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil //todo mark
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

// DoRemoteChannelCheckAndSetupWebhook https://developer.paypal.com/docs/subscriptions/webhooks/
func (p Paypal) DoRemoteChannelCheckAndSetupWebhook(ctx context.Context, payChannel *entity.OverseaPayChannel) (err error) {
	utility.Assert(payChannel != nil, "payChannel is nil")
	client, _ := NewClient(payChannel.ChannelKey, payChannel.ChannelSecret, payChannel.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	result, err := client.ListWebhooks(ctx, paypal.AncorTypeApplication)
	if err != nil {
		return err
	}
	if len(result.Webhooks) > 1 {
		return gerror.New("webhook endpoints count > 1")
	}
	//过滤不可用
	if len(result.Webhooks) == 0 {
		//创建
		_, err := client.CreateWebhook(ctx, &paypal.CreateWebhookRequest{
			URL: out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id)),
			EventTypes: []paypal.WebhookEventType{
				{Name: "BILLING.SUBSCRIPTION.CREATED"},
				{Name: "BILLING.SUBSCRIPTION.ACTIVATED"},
				{Name: "BILLING.SUBSCRIPTION.UPDATED"},
				{Name: "BILLING.SUBSCRIPTION.EXPIRED"},
				{Name: "BILLING.SUBSCRIPTION.CANCELLED"},
				{Name: "BILLING.SUBSCRIPTION.SUSPENDED"},
				{Name: "BILLING.SUBSCRIPTION.PAYMENT.FAILED"},
			},
		})
		if err != nil {
			return err
		}
		if err != nil {
			return nil
		}
		//更新 secret
		//utility.Assert(len(result.Secret) > 0, "secret is nil")
		//err = query.UpdatePayChannelWebhookSecret(ctx, int64(payChannel.Id), result.Secret)
		//if err != nil {
		//	return err
		//}
	} else {
		utility.Assert(len(result.Webhooks) == 1, "internal webhook update, count is not 1")
		//检查并更新, todo mark 优化检查逻辑，如果 evert 一致不用发起更新
		webhook := result.Webhooks[0]
		//utility.Assert(strings.Compare(result.Status, "enabled") == 0, "webhook not status enabled after updated")// todo mark 需要检查里面的每一项
		_, err := client.UpdateWebhook(ctx, webhook.ID, []paypal.WebhookField{
			{
				Operation: "replace",
				Path:      "/event_types",
				Value: []paypal.WebhookEventType{
					{Name: "BILLING.SUBSCRIPTION.CREATED"},
					{Name: "BILLING.SUBSCRIPTION.ACTIVATED"},
					{Name: "BILLING.SUBSCRIPTION.UPDATED"},
					{Name: "BILLING.SUBSCRIPTION.EXPIRED"},
					{Name: "BILLING.SUBSCRIPTION.CANCELLED"},
					{Name: "BILLING.SUBSCRIPTION.SUSPENDED"},
					{Name: "BILLING.SUBSCRIPTION.PAYMENT.FAILED"},
				},
			},
			{
				Operation: "replace",
				Path:      "/url",
				Value:     out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id)),
			},
		})
		if err != nil && strings.Compare(err.(*paypal.ErrorResponse).Name, "WEBHOOK_PATCH_REQUEST_NO_CHANGE") != 0 {
			//WEBHOOK_PATCH_REQUEST_NO_CHANGE 忽略没有更改的错误
			return err
		}
	}

	return nil
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
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	if len(channelEntity.UniqueProductId) > 0 {
		//paypal 保证只创建一个 Product
		return &ro.CreateProductInternalResp{
			ChannelProductId:     channelEntity.UniqueProductId,
			ChannelProductStatus: "",
		}, nil
	}
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
		ImageUrl:    plan.ImageUrl, //paypal 通道可为空
		HomeUrl:     plan.HomeUrl,  //paypal 通道可为空
	})
	if err != nil {
		return nil, err
	}
	err = query.SavePayChannelUniqueProductId(ctx, int64(channelEntity.Id), productResult.ID)
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
	//税费是否包含处理
	taxInclusive := true
	if plan.TaxInclusive == 0 {
		//税费不包含
		taxInclusive = false
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
						Currency: strings.ToUpper(plan.Currency),
						Value:    utility.ConvertFenToYuanMinUnitStr(plan.Amount), //paypal 需要元为单位，小数点处理
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
			Percentage: strconv.Itoa(plan.TaxPercentage),
			Inclusive:  taxInclusive, //传递 false 表示由 paypal 帮助计算税率并加到价格上，true 反之
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

func (p Paypal) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	//TODO implement me

	panic("implement me")
}

func (p Paypal) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
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
