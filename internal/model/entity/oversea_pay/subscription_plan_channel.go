// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlanChannel is the golang structure for table subscription_plan_channel.
type SubscriptionPlanChannel struct {
	Id                   uint64      `json:"id"                   ` //
	GmtCreate            *gtime.Time `json:"gmtCreate"            ` // 创建时间
	GmtModify            *gtime.Time `json:"gmtModify"            ` // 修改时间
	PlanId               int64       `json:"planId"               ` // 计划ID
	ChannelId            int64       `json:"channelId"            ` // 支付渠道Id
	ChannelPlanId        string      `json:"channelPlanId"        ` // 支付渠道plan_Id
	ChannelProductId     string      `json:"channelProductId"     ` // 支付渠道product_Id
	ChannelPlanStatus    string      `json:"channelPlanStatus"    ` // channel_plan_status
	ChannelProductStatus string      `json:"channelProductStatus" ` // channel_product_status
	Data                 string      `json:"data"                 ` // 渠道额外参数，JSON格式
	IsDeleted            int         `json:"isDeleted"            ` //
	Status               int         `json:"status"               ` // 渠道绑定状态，0-Init | 1-Create｜2-Active｜3-Inactive
}
