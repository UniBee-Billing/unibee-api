package api

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/plutov/paypal/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unibee-api/internal/consts"
	"unibee-api/internal/logic/gateway/api/log"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/gateway/util"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

// link：https://developer.paypal.com/docs/api/payments/v1/#payment_create
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

func (p Paypal) GatewayUserPaymentMethodListQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserPaymentMethodListInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserCreate(ctx context.Context, gateway *entity.MerchantGateway, user *entity.UserAccount) (res *ro.GatewayUserCreateInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewaySubscriptionEndTrial(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayPaymentList(ctx context.Context, gateway *entity.MerchantGateway, listReq *ro.GatewayPaymentListReq) (res []*ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundList(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res []*ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayPaymentDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayPaymentId string) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundDetail(ctx context.Context, gateway *entity.MerchantGateway, gatewayRefundId string) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayMerchantBalancesQuery(ctx context.Context, gateway *entity.MerchantGateway) (res *ro.GatewayMerchantBalanceQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayInvoiceCancel(ctx context.Context, gateway *entity.MerchantGateway, cancelInvoiceInternalReq *ro.GatewayCancelInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayUserDetailQuery(ctx context.Context, gateway *entity.MerchantGateway, userId int64) (res *ro.GatewayUserDetailQueryInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayInvoiceCreateAndPay(ctx context.Context, gateway *entity.MerchantGateway, createInvoiceInternalReq *ro.GatewayCreateInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayInvoicePay(ctx context.Context, gateway *entity.MerchantGateway, payInvoiceInternalReq *ro.GatewayPayInvoiceInternalReq) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayInvoiceDetails(ctx context.Context, gateway *entity.MerchantGateway, gatewayInvoiceId string) (res *ro.GatewayDetailInvoiceInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewaySubscriptionNewTrialEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription, newTrialEnd int64) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewaySubscriptionUpdateProrationPreview(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionPreviewInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	// gateway_webhook_entry
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

func (p Paypal) GatewaySubscriptionCreate(ctx context.Context, subscriptionRo *ro.GatewayCreateSubscriptionInternalReq) (res *ro.GatewayCreateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.GatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(subscriptionRo.GatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, subscriptionRo.GatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.GatewayPlan.GatewayPlanId,
		// todo mark
		StartTime:     nil,
		EffectiveTime: nil,
		Quantity:      "",
		//Test Setup Fee
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
	log.SaveChannelHttpLog("GatewaySubscriptionCreate", param, createSubscription, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	//Get Link
	var link string
	for _, item := range createSubscription.Links {
		if strings.Compare(item.Rel, "approve") == 0 {
			link = item.Href
		}
	}
	jsonData, _ := gjson.Marshal(createSubscription)
	return &ro.GatewayCreateSubscriptionInternalResp{
		GatewayUserId:             createSubscription.CustomID,
		Link:                      link,
		GatewaySubscriptionId:     createSubscription.ID,
		GatewaySubscriptionStatus: string(createSubscription.SubscriptionStatus),
		Data:                      string(jsonData),
		Status:                    0, //todo mark
	}, nil
}

func (p Paypal) GatewaySubscriptionCancel(ctx context.Context, subscriptionCancelInternalReq *ro.GatewayCancelSubscriptionInternalReq) (res *ro.GatewayCancelSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

// todo mark paypal 的 cancel 似乎是无法恢复的，和 stripe 不一样，需要确认是否有真实 cancel 的需求
func (p Paypal) GatewaySubscriptionCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelAtPeriodEndSubscriptionInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(gatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	err = client.CancelSubscription(ctx, subscription.GatewaySubscriptionId, "")
	log.SaveChannelHttpLog("GatewaySubscriptionCancelAtPeriodEnd", nil, nil, err, "", nil, gateway)
	if err != nil {
		return nil, err
	} // cancelReason

	return &ro.GatewayCancelAtPeriodEndSubscriptionInternalResp{}, nil //todo mark
}

func (p Paypal) GatewaySubscriptionCancelLastCancelAtPeriodEnd(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayCancelLastCancelAtPeriodEndSubscriptionInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

// Int returns a pointer to the int64 value passed in.
func Int(v int) *int {
	return &v
}

// GatewaySubscriptionUpdate 新旧 Plan 需要在同一个 Product 下，Product Seems Useless，stripe don't
// Need Update After Paid，but stripe don't
func (p Paypal) GatewaySubscriptionUpdate(ctx context.Context, subscriptionRo *ro.GatewayUpdateSubscriptionInternalReq) (res *ro.GatewayUpdateSubscriptionInternalResp, err error) {
	utility.Assert(subscriptionRo.GatewayPlan.GatewayId > 0, "Id Invalid")
	utility.Assert(len(subscriptionRo.GatewayPlan.GatewayProductId) > 0, "Product Not Created")
	gateway := util.GetGatewayById(ctx, subscriptionRo.GatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "Gateway Not Found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	param := paypal.SubscriptionBase{
		PlanID: subscriptionRo.GatewayPlan.GatewayPlanId,
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
				//			Amount:    utility.ConvertCentToDollarStr(subscriptionRo.Plan.Amount * 2), //paypal need doller，cents should change
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
					Value:    "25", //todo mark
				},
				SetupFeeFailureAction:   paypal.SetupFeeFailureActionCancel,
				PaymentFailureThreshold: 2,
			},
			Taxes: nil,
		},
		//todo mark
	}
	updateSubscription, err := client.ReviseSubscription(ctx, subscriptionRo.Subscription.GatewaySubscriptionId, param)
	log.SaveChannelHttpLog("GatewaySubscriptionUpdate", param, updateSubscription, err, subscriptionRo.Subscription.GatewaySubscriptionId, nil, gateway)
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
	return &ro.GatewayUpdateSubscriptionInternalResp{
		GatewayUpdateId: updateSubscription.ID,
		Data:            string(jsonData),
		Link:            link,
		Paid:            false,
	}, nil
}

func (p Paypal) GatewaySubscriptionDetails(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan, subscription *entity.Subscription) (res *ro.GatewayDetailSubscriptionInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(gatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	response, err := client.GetSubscriptionDetails(ctx, subscription.GatewaySubscriptionId)
	log.SaveChannelHttpLog("GatewaySubscriptionDetails", subscription.GatewaySubscriptionId, response, err, "", nil, gateway)
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

	return &ro.GatewayDetailSubscriptionInternalResp{
		Status:        status,
		GatewayStatus: string(response.SubscriptionStatus),
		Data:          utility.FormatToJsonString(response),
	}, nil
}

func (p Paypal) GatewayPlanActive(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(gatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.ActivateSubscriptionPlan(ctx, gatewayPlan.GatewayPlanId)
	log.SaveChannelHttpLog("GatewayPlanActive", gatewayPlan.GatewayPlanId, nil, err, "", nil, gateway)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) GatewayPlanDeactivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(gatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = client.DeactivateSubscriptionPlans(ctx, gatewayPlan.GatewayPlanId)
	log.SaveChannelHttpLog("GatewayPlanDeactivate", gatewayPlan.GatewayPlanId, nil, err, "", nil, gateway)
	if err != nil {
		return err
	}
	return nil
}

func (p Paypal) GatewayProductCreate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreateProductInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	if len(gateway.UniqueProductId) > 0 {
		//paypal 保证只创建一个 Product
		return &ro.GatewayCreateProductInternalResp{
			GatewayProductId:     gateway.UniqueProductId,
			GatewayProductStatus: "",
		}, nil
	}
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
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
	log.SaveChannelHttpLog("GatewayProductCreate", param, productResult, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	err = query.SaveGatewayUniqueProductId(ctx, int64(gateway.Id), productResult.ID)
	if err != nil {
		return nil, err
	}
	return &ro.GatewayCreateProductInternalResp{
		GatewayProductId:     productResult.ID,
		GatewayProductStatus: "",
	}, nil
}

func (p Paypal) GatewayPlanCreateAndActivate(ctx context.Context, plan *entity.SubscriptionPlan, gatewayPlan *entity.GatewayPlan) (res *ro.GatewayCreatePlanInternalResp, err error) {
	utility.Assert(gatewayPlan.GatewayId > 0, "Gateway Not Found")
	utility.Assert(len(gatewayPlan.GatewayProductId) > 0, "Product未创建")
	gateway := util.GetGatewayById(ctx, gatewayPlan.GatewayId)
	utility.Assert(gateway != nil, "gateway not found")
	client, _ := NewClient(gateway.GatewayKey, gateway.GatewaySecret, gateway.Host)
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
		ProductId:   gatewayPlan.GatewayProductId,
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
						Value:    utility.ConvertCentToDollarStr(plan.Amount, plan.Currency),
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
	log.SaveChannelHttpLog("GatewayPlanCreateAndActivate", param, subscriptionPlan, err, "", nil, gateway)
	if err != nil {
		return nil, err
	}
	jsonData, _ := gjson.Marshal(subscriptionPlan)
	return &ro.GatewayCreatePlanInternalResp{
		GatewayPlanId:     subscriptionPlan.ID,
		GatewayPlanStatus: string(subscriptionPlan.Status),
		Data:              string(jsonData),
		Status:            consts.GatewayPlanStatusActive,
	}, nil
}

func (p Paypal) GatewayPayment(ctx context.Context, createPayContext *ro.CreatePayContext) (res *ro.CreatePayInternalResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayCapture(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCaptureRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayCancel(ctx context.Context, payment *entity.Payment) (res *ro.OutPayCancelRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayPayStatusCheck(ctx context.Context, payment *entity.Payment) (res *ro.GatewayPaymentRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefundStatusCheck(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}

func (p Paypal) GatewayRefund(ctx context.Context, payment *entity.Payment, refund *entity.Refund) (res *ro.OutPayRefundRo, err error) {
	//TODO implement me
	panic("implement me")
}
