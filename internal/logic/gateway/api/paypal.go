package api

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/plutov/paypal/v4"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/gateway/api/log"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/gateway/util"
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

func (p Paypal) DoRemoteChannelUserPaymentMethodListQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelUserCreate(ctx context.Context, payChannel *entity.MerchantGateway, user *entity.UserAccount) (res *ro.ChannelUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPaymentList(ctx context.Context, payChannel *entity.MerchantGateway, listReq *ro.ChannelPaymentListReq) (res []*ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefundList(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPaymentDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelPaymentId string) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefundDetail(ctx context.Context, payChannel *entity.MerchantGateway, channelRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelMerchantBalancesQuery(ctx context.Context, payChannel *entity.MerchantGateway) (res *ro.ChannelMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelInvoiceCancel(ctx context.Context, payChannel *entity.MerchantGateway, cancelInvoiceInternalReq *ro.ChannelCancelInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelUserDetailQuery(ctx context.Context, payChannel *entity.MerchantGateway, userId int64) (res *ro.ChannelUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelInvoiceCreateAndPay(ctx context.Context, payChannel *entity.MerchantGateway, createInvoiceInternalReq *ro.ChannelCreateInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelInvoicePay(ctx context.Context, payChannel *entity.MerchantGateway, payInvoiceInternalReq *ro.ChannelPayInvoiceInternalReq) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelInvoiceDetails(ctx context.Context, payChannel *entity.MerchantGateway, channelInvoiceId string) (res *ro.ChannelDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelSubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionPreviewInternalResp, err error) {
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
	utility.Assert(subscriptionRo.PlanChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.GatewayPlanId,
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
							Value:    utility.ConvertCentToDollarStr(subscriptionRo.Subscription.Amount, subscriptionRo.Subscription.Currency), //paypal 需要元为单位，小数点处理
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

func (p Paypal) DoRemoteChannelSubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.ChannelCancelSubscriptionInternalReq) (res *ro.ChannelCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

// todo mark paypal 的 cancel 似乎是无法恢复的，和 stripe 不一样，需要确认是否有真实 cancel 的需求
func (p Paypal) DoRemoteChannelSubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	err = client.CancelSubscription(ctx, subscription.GatewaySubscriptionId, "")
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionCancelAtPeriodEnd", nil, nil, err, "", nil, channelEntity)
	if err != nil {
		return nil, err
	} // cancelReason

	return &ro.ChannelCancelAtPeriodEndSubscriptionInternalResp{}, nil //todo mark
}

func (p Paypal) DoRemoteChannelSubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

// Int returns a pointer to the int64 value passed in.
func Int(v int) *int {
	return &v
}

// DoRemoteChannelSubscriptionUpdate 新旧 Plan 需要在同一个 Product 下，你这个 Product 有什么用，stripe 不需要
// 需要支付之后才能更新，stripe 不需要
func (p Paypal) DoRemoteChannelSubscriptionUpdate(ctx context.Context, subscriptionRo *ro.ChannelUpdateSubscriptionInternalReq) (res *ro.ChannelUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.PlanChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(subscriptionRo.PlanChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, subscriptionRo.PlanChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.PlanChannel.GatewayPlanId,
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
							Value:    utility.ConvertCentToDollarStr(subscriptionRo.Plan.Amount, subscriptionRo.Plan.Currency), //paypal need float
						},
						CreateTime: time.Now(),
						UpdateTime: time.Now(),
					},
					Sequence: Int(1),
				},
				//{
				//	PricingScheme: paypal.PricingScheme{
				//		InvoiceDate: 1,
				//		FixedPrice: paypal.Money{
				//			Currency: strings.ToUpper(subscriptionRo.Plan.Currency),
				//			Amount:    utility.ConvertCentToDollarStr(subscriptionRo.Plan.Amount * 2), //paypal 需要元为单位，小数点处理
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
	updateSubscription, err := client.ReviseSubscription(ctx, subscriptionRo.Subscription.GatewaySubscriptionId, param)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionUpdate", param, updateSubscription, err, subscriptionRo.Subscription.GatewaySubscriptionId, nil, channelEntity)
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
		ChannelUpdateId: updateSubscription.ID,
		Data:            string(jsonData),
		Link:            link,
		Paid:            false,
	}, nil
}

func (p Paypal) DoRemoteChannelSubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.ChannelDetailSubscriptionInternalResp, err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	response, err := client.GetSubscriptionDetails(ctx, subscription.GatewaySubscriptionId)
	log.SaveChannelHttpLog("DoRemoteChannelSubscriptionDetails", subscription.GatewaySubscriptionId, response, err, "", nil, channelEntity)
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

func (p Paypal) DoRemoteChannelPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.ActivateSubscriptionPlan(ctx, planChannel.GatewayPlanId)
	log.SaveChannelHttpLog("DoRemoteChannelPlanActive", planChannel.GatewayPlanId, nil, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) DoRemoteChannelPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.DeactivateSubscriptionPlans(ctx, planChannel.GatewayPlanId)
	log.SaveChannelHttpLog("DoRemoteChannelPlanDeactivate", planChannel.GatewayPlanId, nil, err, "", nil, channelEntity)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) DoRemoteChannelProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreateProductInternalResp, err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	if len(channelEntity.UniqueProductId) > 0 {
		//paypal 保证只创建一个 Product
		return &ro.ChannelCreateProductInternalResp{
			GatewayProductId:     channelEntity.UniqueProductId,
			ChannelProductStatus: "",
		}, nil
	}
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.Product{
		Name:        plan.GatewayProductName,
		Description: plan.GatewayProductDescription,
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
		GatewayProductId:     productResult.ID,
		ChannelProductStatus: "",
	}, nil
}

func (p Paypal) DoRemoteChannelPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, planChannel *entity.GatewayPlan) (res *ro.ChannelCreatePlanInternalResp, err error) {
	utility.Assert(planChannel.GatewayId > 0, "支付渠道异常")
	utility.Assert(len(planChannel.GatewayProductId) > 0, "Product未创建")
	channelEntity := util.GetOverseaPayChannel(ctx, planChannel.GatewayId)
	utility.Assert(channelEntity != nil, "out channel not found")
	client, _ := NewClient(channelEntity.GatewayKey, channelEntity.GatewaySecret, channelEntity.Host)
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
		ProductId:   planChannel.GatewayProductId,
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
						Value:    utility.ConvertCentToDollarStr(plan.Amount, plan.Currency), //paypal 需要元为单位，小数点处理
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
			Percentage: strconv.Itoa(plan.TaxScale), // todo mark
			Inclusive:  taxInclusive,                //传递 false 表示由 paypal 帮助计算税率并加到价格上，true 反之
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
		GatewayPlanId:     subscriptionPlan.ID,
		ChannelPlanStatus: string(subscriptionPlan.Status),
		Data:              string(jsonData),
		Status:            consts.PlanChannelStatusActive,
	}, nil
}

func (p Paypal) DoRemoteChannelPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.ChannelPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) DoRemoteChannelRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
