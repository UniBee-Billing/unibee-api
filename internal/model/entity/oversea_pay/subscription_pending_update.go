// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPendingUpdate is the golang structure for table subscription_pending_update.
type SubscriptionPendingUpdate struct {
	Id                   uint64      `json:"id"                   ` //
	MerchantId           int64       `json:"merchantId"           ` // 商户Id
	SubscriptionId       string      `json:"subscriptionId"       ` // 订阅id（内部编号）
	UpdateSubscriptionId string      `json:"updateSubscriptionId" ` // 升级单ID（内部编号）
	GmtCreate            *gtime.Time `json:"gmtCreate"            ` // 创建时间
	Amount               int64       `json:"amount"               ` // 金额,单位：分
	UpdateAmount         int64       `json:"updateAmount"         ` // 升级到金额,单位：分
	Currency             string      `json:"currency"             ` // 货币
	UpdateCurrency       string      `json:"updateCurrency"       ` // 升级到货币
	PlanId               int64       `json:"planId"               ` // 计划ID
	UpdatePlanId         int64       `json:"updatePlanId"         ` // 升级到计划ID
	Quantity             int64       `json:"quantity"             ` // quantity
	UpdateQuantity       int64       `json:"updateQuantity"       ` // 升级到quantity
	AddonData            string      `json:"addonData"            ` // plan addon json data
	UpdatedAddonData     string      `json:"updatedAddonData"     ` // 升级到plan addon json data
	ChannelId            int64       `json:"channelId"            ` // 支付渠道Id
	Status               int         `json:"status"               ` // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Suspend
	UserId               int64       `json:"userId"               ` // userId
	ChannelUpdateId      string      `json:"channelUpdateId"      ` // 支付渠道订阅更新单id
	Data                 string      `json:"data"                 ` // 渠道额外参数，JSON格式
	ResponseData         string      `json:"responseData"         ` // 渠道返回参数，JSON格式
	GmtModify            *gtime.Time `json:"gmtModify"            ` // 修改时间
	IsDeleted            int         `json:"isDeleted"            ` //
	Link                 string      `json:"link"                 ` //
	ChannelStatus        string      `json:"channelStatus"        ` // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
	ChannelInvoiceId     string      `json:"channelInvoiceId"     ` // 关联渠道发票 Id
}
