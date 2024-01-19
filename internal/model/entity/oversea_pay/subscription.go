// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscription is the golang structure for table subscription.
type Subscription struct {
	Id                     uint64      `json:"id"                     ` //
	SubscriptionId         string      `json:"subscriptionId"         ` // 订阅id（内部编号）
	UserId                 int64       `json:"userId"                 ` // userId
	GmtCreate              *gtime.Time `json:"gmtCreate"              ` // 创建时间
	GmtModify              *gtime.Time `json:"gmtModify"              ` // 修改时间
	Amount                 int64       `json:"amount"                 ` // 金额,单位：分
	Currency               string      `json:"currency"               ` // 货币
	MerchantId             int64       `json:"merchantId"             ` // 商户Id
	PlanId                 int64       `json:"planId"                 ` // 计划ID
	Quantity               int64       `json:"quantity"               ` // quantity
	AddonData              string      `json:"addonData"              ` // plan addon json data
	ChannelId              int64       `json:"channelId"              ` // 支付渠道Id
	Status                 int         `json:"status"                 ` // 订阅单状态，0-Init | 1-Create｜2-Active｜3-PendingInActive | 4-Cancel | 5-Expire | 6- Suspend
	ChannelSubscriptionId  string      `json:"channelSubscriptionId"  ` // 支付渠道订阅id
	Data                   string      `json:"data"                   ` // 渠道额外参数，JSON格式
	ResponseData           string      `json:"responseData"           ` // 渠道返回参数，JSON格式
	ChannelUserId          string      `json:"channelUserId"          ` // 渠道用户 Id
	CustomerName           string      `json:"customerName"           ` // customer_name
	CustomerEmail          string      `json:"customerEmail"          ` // customer_email
	IsDeleted              int         `json:"isDeleted"              ` //
	Link                   string      `json:"link"                   ` //
	ChannelStatus          string      `json:"channelStatus"          ` // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelLatestInvoiceId string      `json:"channelLatestInvoiceId" ` // 渠道最新发票 id
	ChannelItemData        string      `json:"channelItemData"        ` // channel_item_data
	CancelAtPeriodEnd      int         `json:"cancelAtPeriodEnd"      ` // 是否在周期结束时取消，0-false | 1-true
	CurrentPeriodStart     int64       `json:"currentPeriodStart"     ` // current_period_start
	CurrentPeriodEnd       int64       `json:"currentPeriodEnd"       ` // current_period_end
	TrialEnd               int64       `json:"trialEnd"               ` // trial_end
	ReturnUrl              string      `json:"returnUrl"              ` //
	FirstPayTime           *gtime.Time `json:"firstPayTime"           ` // 首次支付时间
	CancelReason           string      `json:"cancelReason"           ` //
	CountryCode            string      `json:"countryCode"            ` //
	VatNumber              string      `json:"vatNumber"              ` //
	TaxPercentage          int64       `json:"taxPercentage"          ` // Tax税率，万分位，1000 表示 10%
	VatVerifyData          string      `json:"vatVerifyData"          ` //
}
