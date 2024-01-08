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
	UpdateSubscriptionId   interface{} // 升级来源订阅 ID（内部编号）
	GmtCreate              *gtime.Time // 创建时间
	Amount                 interface{} // 金额,单位：分
	Currency               interface{} // 货币
	MerchantId             interface{} // 商户Id
	PlanId                 interface{} // 计划ID
	Quantity               interface{} // quantity
	AddonData              interface{} // plan addon json data
	ChannelId              interface{} // 支付渠道Id
	Status                 interface{} // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend | 4-Cancel | 5-Expire
	UserId                 interface{} // userId
	ChannelSubscriptionId  interface{} // 支付渠道订阅id
	Data                   interface{} // 渠道额外参数，JSON格式
	ResponseData           interface{} // 渠道返回参数，JSON格式
	ChannelUserId          interface{} // 渠道用户 Id
	CustomerName           interface{} // customer_name
	CustomerEmail          interface{} // customer_email
	GmtModify              *gtime.Time // 修改时间
	IsDeleted              interface{} //
	Link                   interface{} //
	ChannelStatus          interface{} // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelLatestInvoiceId interface{} // 渠道最新发票 id
	CancelAtPeriodEnd      interface{} // 是否在周期结束时取消，0-false | 1-true
}
