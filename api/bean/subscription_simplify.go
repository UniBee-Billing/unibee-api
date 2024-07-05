package bean

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/oversea_pay"
)

type SubscriptionConfig struct {
	DowngradeEffectImmediately         bool   `json:"downgradeEffectImmediately" dc:"DowngradeEffectImmediately, whether subscription update should effect immediately or at period end, default at period end"`
	UpgradeProration                   bool   `json:"upgradeProration" dc:"UpgradeProration, whether subscription update generation proration invoice or not, default yes"`
	IncompleteExpireTime               int64  `json:"incompleteExpireTime" dc:"IncompleteExpireTime, em.. default 1day for plan of month type"`
	InvoiceEmail                       bool   `json:"invoiceEmail" dc:"InvoiceEmail, whether to send invoice email to user, default yes"`
	InvoicePdfGenerate                 bool   `json:"invoicePdfGenerate" dc:"InvoicePdfGenerate, whether to generate invoice pdf to user, default yes"`
	TryAutomaticPaymentBeforePeriodEnd int64  `json:"tryAutomaticPaymentBeforePeriodEnd" dc:"TryAutomaticPaymentBeforePeriodEnd, default 30 min"`
	GatewayVATRule                     string `json:"gatewayVATRule" dc:""`
}

type SubscriptionSimplify struct {
	Id                     uint64                 `json:"id"                          description:""`                                                                                                                                                               //
	SubscriptionId         string                 `json:"subscriptionId"              description:"subscription id"`                                                                                                                                                // subscription id
	UserId                 uint64                 `json:"userId"                      description:"userId"`                                                                                                                                                         // userId
	TaskTime               *gtime.Time            `json:"taskTime"                    description:"task_time"`                                                                                                                                                      // task_time
	Amount                 int64                  `json:"amount"                      description:"amount, cent"`                                                                                                                                                   // amount, cent
	Currency               string                 `json:"currency"                    description:"currency"`                                                                                                                                                       // currency
	MerchantId             uint64                 `json:"merchantId"                  description:"merchant id"`                                                                                                                                                    // merchant id
	PlanId                 uint64                 `json:"planId"                      description:"plan id"`                                                                                                                                                        // plan id
	Quantity               int64                  `json:"quantity"                    description:"quantity"`                                                                                                                                                       // quantity
	AddonData              string                 `json:"addonData"                   description:"plan addon json data"`                                                                                                                                           // plan addon json data
	LatestInvoiceId        string                 `json:"latestInvoiceId"             description:"latest_invoice_id"`                                                                                                                                              // latest_invoice_id
	Type                   int                    `json:"type"                        description:"sub type, 0-gateway sub, 1-unibee sub"`                                                                                                                          // sub type, 0-gateway sub, 1-unibee sub
	GatewayId              uint64                 `json:"gatewayId"                   description:"gateway_id"`                                                                                                                                                     // gateway_id
	Status                 int                    `json:"status"                      description:"status，1-Pending｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete | 8-Processing"`                                                    // status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	Link                   string                 `json:"link"                        description:""`                                                                                                                                                               //
	GatewayStatus          string                 `json:"gatewayStatus"               description:"gateway status，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get"` // gateway status，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	Features               string                 `json:"features"                    description:"features"`                                                                                                                                                       // gateway_item_data
	CancelAtPeriodEnd      int                    `json:"cancelAtPeriodEnd"           description:"whether cancel at period end，0-false | 1-true"`                                                                                                                  // whether cancel at period end，0-false | 1-true
	LastUpdateTime         int64                  `json:"lastUpdateTime"              description:""`                                                                                                                                                               //
	CurrentPeriodStart     int64                  `json:"currentPeriodStart"          description:"current_period_start, utc time"`                                                                                                                                 // current_period_start, utc time
	CurrentPeriodEnd       int64                  `json:"currentPeriodEnd"            description:"current_period_end, utc time"`                                                                                                                                   // current_period_end, utc time
	BillingCycleAnchor     int64                  `json:"billingCycleAnchor"          description:"billing_cycle_anchor"`                                                                                                                                           // billing_cycle_anchor
	DunningTime            int64                  `json:"dunningTime"                 description:"dunning_time, utc time"`                                                                                                                                         // dunning_time, utc time
	TrialEnd               int64                  `json:"trialEnd"                    description:"trial_end, utc time"`                                                                                                                                            // trial_end, utc time
	ReturnUrl              string                 `json:"returnUrl"                   description:""`                                                                                                                                                               //
	FirstPaidTime          int64                  `json:"firstPaidTime"               description:"first success payment time"`                                                                                                                                     // first success payment time
	CancelReason           string                 `json:"cancelReason"                description:""`                                                                                                                                                               //
	CountryCode            string                 `json:"countryCode"                 description:""`                                                                                                                                                               //
	VatNumber              string                 `json:"vatNumber"                   description:""`                                                                                                                                                               //
	TaxPercentage          int64                  `json:"taxPercentage"               description:"TaxPercentage，1000 = 10%"`                                                                                                                                       // Tax Percentage，1000 = 10%
	PendingUpdateId        string                 `json:"pendingUpdateId"             description:""`                                                                                                                                                               //
	CreateTime             int64                  `json:"createTime"                  description:"create utc time"`                                                                                                                                                // create utc time
	TestClock              int64                  `json:"testClock"                   description:"test_clock, simulator clock for subscription, if set , sub will out of cronjob controll"`                                                                        // test_clock, simulator clock for subscription, if set , sub will out of cronjob controll
	Metadata               map[string]interface{} `json:"metadata" description:""`
	GasPayer               string                 `json:"gasPayer"                  description:"who pay the gas, merchant|user"` // who pay the gas, merchant|user
	DefaultPaymentMethodId string                 `json:"defaultPaymentMethodId"    description:""`
}

func SimplifySubscription(one *entity.Subscription) *SubscriptionSimplify {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &SubscriptionSimplify{
		Id:                     one.Id,
		SubscriptionId:         one.SubscriptionId,
		UserId:                 one.UserId,
		TaskTime:               one.TaskTime,
		Amount:                 one.Amount,
		Currency:               one.Currency,
		MerchantId:             one.MerchantId,
		PlanId:                 one.PlanId,
		Quantity:               one.Quantity,
		AddonData:              one.AddonData,
		LatestInvoiceId:        one.LatestInvoiceId,
		Type:                   one.Type,
		GatewayId:              one.GatewayId,
		Status:                 one.Status,
		Link:                   one.Link,
		GatewayStatus:          one.GatewayStatus,
		Features:               one.GatewayItemData,
		CancelAtPeriodEnd:      one.CancelAtPeriodEnd,
		LastUpdateTime:         one.LastUpdateTime,
		CurrentPeriodStart:     one.CurrentPeriodStart,
		CurrentPeriodEnd:       one.CurrentPeriodEnd,
		BillingCycleAnchor:     one.BillingCycleAnchor,
		DunningTime:            one.DunningTime,
		TrialEnd:               one.TrialEnd,
		ReturnUrl:              one.ReturnUrl,
		FirstPaidTime:          one.FirstPaidTime,
		CancelReason:           one.CancelReason,
		CountryCode:            one.CountryCode,
		VatNumber:              one.VatNumber,
		TaxPercentage:          one.TaxPercentage,
		PendingUpdateId:        one.PendingUpdateId,
		CreateTime:             one.CreateTime,
		TestClock:              one.TestClock,
		Metadata:               metadata,
		GasPayer:               one.GasPayer,
		DefaultPaymentMethodId: one.GatewayDefaultPaymentMethod,
	}
}

type SubscriptionOnetimeAddonSimplify struct {
	Id             uint64                 `json:"id"             description:"id"`                                            // id
	SubscriptionId string                 `json:"subscriptionId" description:"subscription_id"`                               // subscription_id
	AddonId        uint64                 `json:"addonId"        description:"onetime addonId"`                               // onetime addonId
	Quantity       int64                  `json:"quantity"       description:"quantity"`                                      // quantity
	Status         int                    `json:"status"         description:"status, 1-create, 2-paid, 3-cancel, 4-expired"` // status, 1-create, 2-paid, 3-cancel, 4-expired
	IsDeleted      int                    `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"`                         // 0-UnDeleted，1-Deleted
	CreateTime     int64                  `json:"createTime"     description:"create utc time"`                               // create utc time
	PaymentId      string                 `json:"paymentId"     description:"PaymentId"`                                      // PaymentId
	PaymentLink    string                 `json:"paymentLink"     description:"PaymentLink"`                                  // PaymentLink
	Metadata       map[string]interface{} `json:"metadata" description:""`
}

func SimplifySubscriptionOnetimeAddonSimplify(one *entity.SubscriptionOnetimeAddon) *SubscriptionOnetimeAddonSimplify {
	if one == nil {
		return nil
	}
	var metadata = make(map[string]interface{})
	if len(one.MetaData) > 0 {
		err := gjson.Unmarshal([]byte(one.MetaData), &metadata)
		if err != nil {
			fmt.Printf("SimplifySubscription Unmarshal Metadata error:%s", err.Error())
		}
	}
	return &SubscriptionOnetimeAddonSimplify{
		Id:             one.Id,
		SubscriptionId: one.SubscriptionId,
		AddonId:        one.AddonId,
		Quantity:       one.Quantity,
		Status:         one.Status,
		IsDeleted:      one.IsDeleted,
		CreateTime:     one.CreateTime,
		PaymentId:      one.PaymentId,
		PaymentLink:    link.GetPaymentLink(one.PaymentId),
		Metadata:       metadata,
	}
}
