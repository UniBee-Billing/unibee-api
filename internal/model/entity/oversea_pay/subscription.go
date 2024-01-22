// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscription is the golang structure for table subscription.
type Subscription struct {
	Id                     uint64      `json:"id"                     description:""`                                                                                                                                                       //
	SubscriptionId         string      `json:"subscriptionId"         description:"订阅id（内部编号）"`                                                                                                                                             // 订阅id（内部编号）
	UserId                 int64       `json:"userId"                 description:"userId"`                                                                                                                                                 // userId
	GmtCreate              *gtime.Time `json:"gmtCreate"              description:"创建时间"`                                                                                                                                                   // 创建时间
	GmtModify              *gtime.Time `json:"gmtModify"              description:"修改时间"`                                                                                                                                                   // 修改时间
	Amount                 int64       `json:"amount"                 description:"金额,单位：分"`                                                                                                                                                // 金额,单位：分
	Currency               string      `json:"currency"               description:"货币"`                                                                                                                                                     // 货币
	MerchantId             int64       `json:"merchantId"             description:"商户Id"`                                                                                                                                                   // 商户Id
	PlanId                 int64       `json:"planId"                 description:"计划ID"`                                                                                                                                                   // 计划ID
	Quantity               int64       `json:"quantity"               description:"quantity"`                                                                                                                                               // quantity
	AddonData              string      `json:"addonData"              description:"plan addon json data"`                                                                                                                                   // plan addon json data
	ChannelId              int64       `json:"channelId"              description:"支付渠道Id"`                                                                                                                                                 // 支付渠道Id
	Status                 int         `json:"status"                 description:"订阅单状态，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend"`                                                                  // 订阅单状态，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend
	ChannelSubscriptionId  string      `json:"channelSubscriptionId"  description:"支付渠道订阅id"`                                                                                                                                               // 支付渠道订阅id
	ChannelUserId          string      `json:"channelUserId"          description:"渠道用户 Id"`                                                                                                                                                // 渠道用户 Id
	CustomerName           string      `json:"customerName"           description:"customer_name"`                                                                                                                                          // customer_name
	CustomerEmail          string      `json:"customerEmail"          description:"customer_email"`                                                                                                                                         // customer_email
	IsDeleted              int         `json:"isDeleted"              description:""`                                                                                                                                                       //
	Link                   string      `json:"link"                   description:""`                                                                                                                                                       //
	ChannelStatus          string      `json:"channelStatus"          description:"渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get"` // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelItemData        string      `json:"channelItemData"        description:"channel_item_data"`                                                                                                                                      // channel_item_data
	CancelAtPeriodEnd      int         `json:"cancelAtPeriodEnd"      description:"是否在周期结束时取消，0-false | 1-true"`                                                                                                                            // 是否在周期结束时取消，0-false | 1-true
	ChannelLatestInvoiceId string      `json:"channelLatestInvoiceId" description:"渠道最新发票 id"`                                                                                                                                              // 渠道最新发票 id
	CurrentPeriodStart     int64       `json:"currentPeriodStart"     description:"current_period_start"`                                                                                                                                   // current_period_start
	CurrentPeriodEnd       int64       `json:"currentPeriodEnd"       description:"current_period_end"`                                                                                                                                     // current_period_end
	CurrentPeriodStartTime *gtime.Time `json:"currentPeriodStartTime" description:""`                                                                                                                                                       //
	CurrentPeriodEndTime   *gtime.Time `json:"currentPeriodEndTime"   description:""`                                                                                                                                                       //
	BillingCycleAnchor     int64       `json:"billingCycleAnchor"     description:"billing_cycle_anchor"`                                                                                                                                   // billing_cycle_anchor
	TrialEnd               int64       `json:"trialEnd"               description:"trial_end"`                                                                                                                                              // trial_end
	ReturnUrl              string      `json:"returnUrl"              description:""`                                                                                                                                                       //
	FirstPayTime           *gtime.Time `json:"firstPayTime"           description:"首次支付时间"`                                                                                                                                                 // 首次支付时间
	CancelReason           string      `json:"cancelReason"           description:""`                                                                                                                                                       //
	CountryCode            string      `json:"countryCode"            description:""`                                                                                                                                                       //
	VatNumber              string      `json:"vatNumber"              description:""`                                                                                                                                                       //
	TaxPercentage          int64       `json:"taxPercentage"          description:"Tax税率，万分位，1000 表示 10%"`                                                                                                                                  // Tax税率，万分位，1000 表示 10%
	VatVerifyData          string      `json:"vatVerifyData"          description:""`                                                                                                                                                       //
	Data                   string      `json:"data"                   description:"渠道额外参数，JSON格式"`                                                                                                                                          // 渠道额外参数，JSON格式
	ResponseData           string      `json:"responseData"           description:"渠道返回参数，JSON格式"`                                                                                                                                          // 渠道返回参数，JSON格式
}
