// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Subscription is the golang structure for table subscription.
type Subscription struct {
	Id                    uint64      `json:"id"                    ` //
	SubscriptionId        string      `json:"subscriptionId"        ` // 内部订阅id
	GmtCreate             *gtime.Time `json:"gmtCreate"             ` // 创建时间
	Type                  int         `json:"type"                  ` // 类型，1-普通单，2-升级单
	UpdateFromId          int64       `json:"updateFromId"          ` // 升级来源 ID
	UpdateToId            int64       `json:"updateToId"            ` // 升级去向 ID
	CompanyId             int64       `json:"companyId"             ` // 公司ID
	MerchantId            int64       `json:"merchantId"            ` // 商户Id
	PlanId                int64       `json:"planId"                ` // 计划ID
	ChannelId             int64       `json:"channelId"             ` // 支付渠道Id
	UserId                int64       `json:"userId"                ` // userId
	Quantity              int64       `json:"quantity"              ` // quantity
	ChannelSubscriptionId string      `json:"channelSubscriptionId" ` // 支付渠道订阅id
	Data                  string      `json:"data"                  ` // 渠道额外参数，JSON格式
	ResponseData          string      `json:"responseData"          ` // 渠道返回参数，JSON格式
	Status                int         `json:"status"                ` // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelUserId         string      `json:"channelUserId"         ` // 渠道用户 Id
	CustomerName          string      `json:"customerName"          ` // customer_name
	CustomerEmail         string      `json:"customerEmail"         ` // customer_email
	GmtModify             *gtime.Time `json:"gmtModify"             ` // 修改时间
	IsDeleted             int         `json:"isDeleted"             ` //
	Link                  string      `json:"link"                  ` //
	ChannelStatus         string      `json:"channelStatus"         ` // 渠道最新状态，Stripe：https://stripe.com/docs/billing/subscriptions/webhooks  Paypal：https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_get
}
