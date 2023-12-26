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
	GmtCreate             *gtime.Time `json:"gmtCreate"             ` // 创建时间
	GmtModify             *gtime.Time `json:"gmtModify"             ` // 修改时间
	CompanyId             int64       `json:"companyId"             ` // 公司ID
	MerchantId            int64       `json:"merchantId"            ` // 商户Id
	PlanId                int64       `json:"planId"                ` // 计划ID
	ChannelId             int64       `json:"channelId"             ` // 支付渠道Id
	UserId                int64       `json:"userId"                ` // userId
	Quantity              int64       `json:"quantity"              ` // quantity
	SubscriptionId        string      `json:"subscriptionId"        ` // 内部订阅id
	ChannelSubscriptionId string      `json:"channelSubscriptionId" ` // 支付渠道订阅id
	Data                  string      `json:"data"                  ` // 渠道额外参数，JSON格式
	ResponseData          string      `json:"responseData"          ` // 渠道返回参数，JSON格式
	IsDeleted             int         `json:"isDeleted"             ` //
	Status                int         `json:"status"                ` // 订阅单状态，0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelUserId         string      `json:"channelUserId"         ` // 渠道用户 Id
	CustomerName          string      `json:"customerName"          ` // customer_name
	CustomerEmail         string      `json:"customerEmail"         ` // customer_email
	Link                  string      `json:"link"                  ` //
}
