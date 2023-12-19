// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlanChannel is the golang structure for table subscription_plan_channel.
type SubscriptionPlanChannel struct {
	Id               uint64      `json:"id"               ` //
	GmtCreate        *gtime.Time `json:"gmtCreate"        ` // 创建时间
	GmtModify        *gtime.Time `json:"gmtModify"        ` // 修改时间
	PlanId           int64       `json:"planId"           ` // 计划ID
	ChannelId        int64       `json:"channelId"        ` // 支付渠道Id
	ChannelPlanId    int64       `json:"channelPlanId"    ` // 支付渠道plan_Id
	ChannelProductId int64       `json:"channelProductId" ` // 支付渠道product_Id
	Data             string      `json:"data"             ` // 渠道额外参数，JSON格式
	IsDeleted        int         `json:"isDeleted"        ` //
}
