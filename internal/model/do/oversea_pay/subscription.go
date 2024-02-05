// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscription is the golang structure of table subscription for DAO operations like Where/Data.
type Subscription struct {
	g.Meta                      `orm:"table:subscription, do:true"`
	Id                          interface{} //
	SubscriptionId              interface{} // subscription id
	UserId                      interface{} // userId
	GmtCreate                   *gtime.Time // create time
	GmtModify                   *gtime.Time // update time
	Amount                      interface{} // amount, cent
	Currency                    interface{} // currency
	MerchantId                  interface{} // merchant id
	PlanId                      interface{} // plan id
	Quantity                    interface{} // quantity
	AddonData                   interface{} // plan addon json data
	TaskTime                    *gtime.Time // task_time
	LatestInvoiceId             interface{} // latest_invoice_id
	Type                        interface{} // sub type, 0-gateway sub, 1-unibee sub
	GatewayId                   interface{} // gateway_id
	Status                      interface{} // status，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	GatewaySubscriptionId       interface{} // gateway subscription id
	CustomerName                interface{} // customer_name
	CustomerEmail               interface{} // customer_email
	IsDeleted                   interface{} // 0-UnDeleted，1-Deleted
	GatewayDefaultPaymentMethod interface{} //
	Link                        interface{} //
	GatewayStatus               interface{} // gateway status，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	GatewayItemData             interface{} // gateway_item_data
	CancelAtPeriodEnd           interface{} // whether cancel at period end，0-false | 1-true
	GatewayLatestInvoiceId      interface{} // gateway latest invoice id
	LastUpdateTime              interface{} //
	CurrentPeriodStart          interface{} // current_period_start
	CurrentPeriodEnd            interface{} // current_period_end
	CurrentPeriodStartTime      *gtime.Time //
	CurrentPeriodEndTime        *gtime.Time //
	BillingCycleAnchor          interface{} // billing_cycle_anchor
	DunningTime                 interface{} //
	TrialEnd                    interface{} // trial_end
	ReturnUrl                   interface{} //
	FirstPayTime                *gtime.Time // first success payment time
	CancelReason                interface{} //
	CountryCode                 interface{} //
	VatNumber                   interface{} //
	TaxScale                    interface{} // Tax Scale，1000 = 10%
	VatVerifyData               interface{} //
	Data                        interface{} //
	ResponseData                interface{} //
	PendingUpdateId             interface{} //
}
