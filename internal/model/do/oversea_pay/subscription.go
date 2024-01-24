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
	g.Meta                 `orm:"table:subscription, do:true"`
	Id                     interface{} //
	SubscriptionId         interface{} // 订阅id（内部编号）
	UserId                 interface{} // userId
	GmtCreate              *gtime.Time // 创建时间
	GmtModify              *gtime.Time // 修改时间
	Amount                 interface{} // 金额,单位：分
	Currency               interface{} // 货币
	MerchantId             interface{} // 商户Id
	PlanId                 interface{} // 计划ID
	Quantity               interface{} // quantity
	AddonData              interface{} // plan addon json data
	ChannelId              interface{} // 支付渠道Id
	Status                 interface{} // 订阅单状态，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend| 7-Incomplete
	ChannelSubscriptionId  interface{} // 支付渠道订阅id
	ChannelUserId          interface{} // 渠道用户 Id
	CustomerName           interface{} // customer_name
	CustomerEmail          interface{} // customer_email
	IsDeleted              interface{} //
	Link                   interface{} //
	ChannelStatus          interface{} // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelItemData        interface{} // channel_item_data
	CancelAtPeriodEnd      interface{} // 是否在周期结束时取消，0-false | 1-true
	ChannelLatestInvoiceId interface{} // 渠道最新发票 id
	LastUpdateTime         interface{} //
	CurrentPeriodStart     interface{} // current_period_start
	CurrentPeriodEnd       interface{} // current_period_end
	CurrentPeriodStartTime *gtime.Time //
	CurrentPeriodEndTime   *gtime.Time //
	BillingCycleAnchor     interface{} // billing_cycle_anchor
	TrialEnd               interface{} // trial_end
	ReturnUrl              interface{} //
	FirstPayTime           *gtime.Time // 首次支付时间
	CancelReason           interface{} //
	CountryCode            interface{} //
	VatNumber              interface{} //
	TaxPercentage          interface{} // Tax税率，万分位，1000 表示 10%
	VatVerifyData          interface{} //
	Data                   interface{} // 渠道额外参数，JSON格式
	ResponseData           interface{} // 渠道返回参数，JSON格式
	PendingUpdateId        interface{} //
}
