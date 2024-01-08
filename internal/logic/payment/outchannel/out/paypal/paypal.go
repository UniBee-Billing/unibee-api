package paypal

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/plutov/paypal/v4"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel/out"
	"go-oversea-pay/internal/logic/payment/outchannel/out/log"
	"go-oversea-pay/internal/logic/payment/outchannel/ro"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	"go-oversea-pay/internal/logic/subscription/handler"
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

func (p Paypal) DoRemoteChannelSubscriptionUpdatePreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	//注册 channel_webhook_entry
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

func (p Paypal) DoRemoteChannelSubscriptionCreate(ctx context.Context, subscriptionRo *ro.ChannelCreateSubscriptionInternalReq) (res *ro.ChannelCreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.ChannelPlanId,
		// todo mark
		StartTime:     nil,
		EffectiveTime: nil,
		Quantity:      "",
		//测试安装费
		ShippingAmount: &paypal.Money{
			Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
			Value:    "10",
		},
		Plan: &paypal.PlanOverride{
			BillingCycles: []paypal.BillingCycleOverride{
				{
					PricingScheme: paypal.PricingScheme{
						Version: 1,
						FixedPrice: paypal.Money{
							Currency: strings.ToUpper(subscriptionRo.Subscription.Currency),
							Value:    utility.ConvertFenToYuanMinUnitStr(subscriptionRo.Subscription.Amount), //paypal 需要元为单位，小数点处理
						},
						CreateTime: time.Now(),
						UpdateTime: time.Now(),
					},
					Sequence: Int(1),
				},
			},
			PaymentPreferences: &paypal.PaymentPreferencesOverride{
				AutoBillOutstanding: false,
				SetupFee: paypal.Money{
					Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
					Value:    "0",
				},
				SetupFeeFailureAction:   paypal.SetupFeeFailureActionCancel,
				PaymentFailureThreshold: 2,
			},
			Taxes: nil,
		},
		Subscriber:         nil,
		AutoRenewal:        false,
		ApplicationContext: nil,
		CustomID:           "",
	}
	createSubscription, err := client.CreateSubscription(ctx, param)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCreate", param, createSubscription, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	//获取 Link
	var link string
	for _, item := range createSubscription.Links {
		if strings.Compare(item.Rel, "approve") == 0 {
			link = item.Href
		}
	}
	jsonData, _ := gjson.Marshal(createSubscription)
	return &ro.ChannelCreateSubscriptionInternalResp{
		ChannelUserId:             createSubscription.CustomID,
		Link:                      link,
		ChannelSubscriptionId:     createSubscription.ID,
		ChannelSubscriptionStatus: string(createSubscription.SubscriptionStatus),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil
}

// todo mark paypal 的 cancel 似乎是无法恢复的，和 stripe 不一样，需要确认是否有真实 cancel 的需求
func (p Paypal) DoRemoteChannelSubscriptionCancel(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
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
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancel", nil, nil, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	} // cancelReason

	return &ro.ChannelCancelSubscriptionInternalResp{}, nil //todo mark
}

// Int returns a pointer to the int64 value passed in.
func Int(v int) *int {
	return &v
}

// DoRemoteChannelSubscriptionUpdate 新旧 Plan 需要在同一个 Product 下，你这个 Product 有什么用，stripe 不需要
// 需要支付之后才能更新，stripe 不需要
func (p Paypal) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.ChannelPlanId,
		//测试安装费
		ShippingAmount: &paypal.Money{
			Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
			Value:    "15",
		},
		Plan: &paypal.PlanOverride{
			BillingCycles: []paypal.BillingCycleOverride{
				{
					PricingScheme: paypal.PricingScheme{
						Version: 1,
						FixedPrice: paypal.Money{
							Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
							Value:    utility.ConvertFenToYuanMinUnitStr(subscriptionRo.Plan.Amount), //paypal 需要元为单位，小数点处理
						},
						CreateTime: time.Now(),
						UpdateTime: time.Now(),
					},
					Sequence: Int(1),
				},
				//{
				//	PricingScheme: paypal.PricingScheme{
				//		Version: 1,
				//		FixedPrice: paypal.Money{
				//			Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
				//			Value:    utility.ConvertFenToYuanMinUnitStr(subscriptionRo.Plan.Amount * 2), //paypal 需要元为单位，小数点处理
				//		},
				//		CreateTime: time.Now(),
				//		UpdateTime: time.Now(),
				//	},
				//	Sequence: Int(1),
				//},
			},
			PaymentPreferences: &paypal.PaymentPreferencesOverride{
				AutoBillOutstanding: false,
				SetupFee: paypal.Money{
					Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
					Value:    "25", //todo mark 开户费在更新的时候似乎没有用处
				},
				SetupFeeFailureAction:   paypal.SetupFeeFailureActionCancel,
				PaymentFailureThreshold: 2,
			},
			Taxes: nil,
		},
		//todo mark
	}
	updateSubscription, err := client.ReviseSubscription(ctx, subscriptionRo.Subscription.ChannelSubscriptionId, param)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", param, updateSubscription, err, subscriptionRo.Subscription.ChannelSubscriptionId, nil, channelEntity)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(updateSubscription)
	var link string
	for _, item := range updateSubscription.Links {
		if strings.Compare(item.Rel, "approve") == 0 {
			link = item.Href
		}
	}
	return &ro.ChannelUpdateSubscriptionInternalResp{
		ChannelSubscriptionId:     updateSubscription.ID,
		ChannelSubscriptionStatus: string(updateSubscription.SubscriptionStatus),
		Data:                      string(jsonData),
		Link:                      link,
		Status:                    0, //todo mark
	}, nil //todo mark
}

func (p Paypal) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.ChannelProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	response, err := client.GetSubscriptionDetails(ctx, subscription.ChannelSubscriptionId)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionDetails", subscription.ChannelSubscriptionId, response, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}

	var status consts.SubscriptionStatusEnum = consts.SubStatusSuspended
	if strings.Compare(string(response.SubscriptionStatus), "ACTIVE") == 0 {
		status = consts.SubStatusActive
	} else if strings.Compare(string(response.SubscriptionStatus), "APPROVAL_PENDING") == 0 ||
		strings.Compare(string(response.SubscriptionStatus), "APPROVED") == 0 {
		status = consts.SubStatusCreate
	} else if strings.Compare(string(response.SubscriptionStatus), "SUSPENDED") == 0 {
		status = consts.SubStatusSuspended
	} else if strings.Compare(string(response.SubscriptionStatus), "CANCELLED") == 0 {
		status = consts.SubStatusCancelled
	} else if strings.Compare(string(response.SubscriptionStatus), "EXPIRED") == 0 {
		status = consts.SubStatusExpired
	}

	return &ro.ChannelDetailSubscriptionInternalResp{
		Status:        status,
		ChannelStatus: string(response.SubscriptionStatus),
		Data:          utility.FormatToJsonString(response),
	}, nil
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
		param := &paypal.CreateWebhookRequest{
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
		}
		response, err := client.CreateWebhook(ctx, param)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", param, response, err, "", nil, payChannel)
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
		param := []paypal.WebhookField{
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
				Value:     strings.Replace(out.GetPaymentWebhookEntranceUrl(int64(payChannel.Id)), "http://", "https://", 1), //paypal 只支持 https
			},
		}
		response, err := client.UpdateWebhook(ctx, webhook.ID, param)
		log.SaveChannelHttpLog("DoRemoteChannelCheckAndSetupWebhook", param, response, err, webhook.ID, nil, payChannel)
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
	log.SaveChannelHttpLog("DoRemoteChannelPlanActive", planChannel.ChannelPlanId, nil, err, "", nil, channelEntity)
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
	log.SaveChannelHttpLog("DoRemoteChannelPlanDeactivate", planChannel.ChannelPlanId, nil, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreateProductInternalResp, err error) {
	utility.Assert(planChannel.ChannelId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.ChannelId)
	utility.Assert(channelEntity != nil, "支付渠道异常 out channel not found")
	if len(channelEntity.UniqueProductId) > 0 {
		//paypal 保证只创建一个 Product
		return &ro.ChannelCreateProductInternalResp{
			ChannelProductId:     channelEntity.UniqueProductId,
			ChannelProductStatus: "",
		}, nil
	}
	client, _ := NewClient(channelEntity.ChannelKey, channelEntity.ChannelSecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.Product{
		Name:        plan.ChannelProductName,
		Description: plan.ChannelProductDescription,
		Category:    paypal.ProductCategorySoftware,
		Type:        paypal.ProductTypeService,
		ImageUrl:    plan.ImageUrl, //paypal 通道可为空
		HomeUrl:     plan.HomeUrl,  //paypal 通道可为空
	}
	productResult, err := client.CreateProduct(ctx, param)
	log.SaveChannelHttpLog("DoRemoteChannelProductCreate", param, productResult, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	err = query.SavePayChannelUniqueProductId(ctx, int64(channelEntity.Id), productResult.ID)
	if err != nil {
		return nil, err
	}
	return &ro.ChannelCreateProductInternalResp{
		ChannelProductId:     productResult.ID,
		ChannelProductStatus: "",
	}, nil
}

func (p Paypal) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.SubscriptionPlanChannel) (res *ro.ChannelCreatePlanInternalResp, err error) {
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
	param := paypal.SubscriptionPlan{
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
					IntervalUnit:  paypal.IntervalUnit(strings.ToUpper(plan.IntervalUnit)),
					IntervalCount: plan.IntervalCount,
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
	}
	subscriptionPlan, err := client.CreateSubscriptionPlan(ctx, param)
	log.SaveChannelHttpLog("DoRemoteChannelPlanCreateAndActivate", param, subscriptionPlan, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(subscriptionPlan)
	return &ro.ChannelCreatePlanInternalResp{
		ChannelPlanId:     subscriptionPlan.ID,
		ChannelPlanStatus: string(subscriptionPlan.Status),
		Data:              string(jsonData),
		Status:            consts.PlanChannelStatusActive,
	}, nil
}

func (p Paypal) processWebhook(ctx context.Context, eventType string, resource *gjson.Json) error {
	unibSub := query.GetSubscriptionByChannelSubscriptionId(ctx, resource.Get("id").String())
	if unibSub != nil {
		plan := query.GetPlanById(ctx, unibSub.PlanId)
		planChannel := query.GetPlanChannel(ctx, unibSub.PlanId, unibSub.ChannelId)
		details, err := p.DoRemoteChannelSubscriptionDetails(ctx, plan, planChannel, unibSub)
		if err != nil {
			return err
		}

		err = handler.HandleSubscriptionWebhookEvent(ctx, unibSub, eventType, details)
		if err != nil {
			return err
		}
		return nil
	} else {
		return gerror.New("subscription not found on channelSubId:" + resource.Get("id").String())
	}
}

func (p Paypal) DoRemoteChannelWebhook(r *ghttp.Request, payChannel *entity.OverseaPayChannel) {
	jsonData, err := r.GetJson()
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook Get Json failed. %v\n", payChannel.Channel, err)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	client, _ := NewClient(payChannel.ChannelKey, payChannel.ChannelSecret, payChannel.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	signature, err := client.VerifyWebhookSignature(r.Context(), r.Request, jsonData.Get("id").String())
	if err != nil {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification success\n", payChannel.Channel)
		r.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.Compare(signature.VerificationStatus, "SUCCESS") == 0 {
		g.Log().Info(r.Context(), "Receive_Webhook_Channel:", payChannel.Channel, " hook:", jsonData.String())
		eventType := jsonData.Get("event_type").String()
		var responseBack = http.StatusOK
		switch eventType {
		case "BILLING.SUBSCRIPTION.EXPIRED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription deleted for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the deleted subscription.
				// handleSubscriptionCanceled(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err)
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		case "BILLING.SUBSCRIPTION.UPDATED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription updated for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the successful attachment of a PaymentMethod.
				// handleSubscriptionUpdated(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err)
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		case "BILLING.SUBSCRIPTION.CREATED":
			resource := jsonData.GetJson("resource")
			if resource == nil || !resource.Contains("id") {
				g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error parsing webhook resource is nil\n", payChannel.Channel)
				r.Response.WriteHeader(http.StatusBadRequest)
				responseBack = http.StatusBadRequest
			} else {
				g.Log().Infof(r.Context(), "Webhook Channel:%s, Subscription created for %d.", payChannel.Channel, resource.Get("id").String())
				// Then define and call a func to handle the successful attachment of a PaymentMethod.
				// handleSubscriptionCreated(subscription)
				err := p.processWebhook(r.Context(), eventType, resource)
				if err != nil {
					g.Log().Errorf(r.Context(), "Webhook Channel:%s, Error HandleSubscriptionWebhookEvent: %v\n", payChannel.Channel, err)
					r.Response.WriteHeader(http.StatusBadRequest)
					responseBack = http.StatusBadRequest
				}
			}
		default:
			g.Log().Errorf(r.Context(), "Webhook Channel:%s, Unhandled event type: %s\n", payChannel.Channel, eventType)
		}
		r.Response.WriteHeader(http.StatusOK)
		log.SaveChannelHttpLog("DoRemoteChannelWebhook", jsonData, responseBack, err, "", nil, payChannel)
		return
	} else {
		g.Log().Errorf(r.Context(), "⚠️  Webhook Channel:%s, Webhook signature verification failed. %v\n", payChannel.Channel)
		r.Response.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
}

func (p Paypal) DoRemoteChannelRedirect(r *ghttp.Request, payChannel *entity.OverseaPayChannel) (res *ro.ChannelRedirectInternalResp, err error) {
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
