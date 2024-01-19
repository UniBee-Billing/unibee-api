// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlanChannel is the golang structure for table subscription_plan_channel.
type SubscriptionPlanChannel struct {
	Id                   uint64      `json:"id"                   description:""`                                             //
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"创建时间"`                                         // 创建时间
	GmtModify            *gtime.Time `json:"gmtModify"            description:"修改时间"`                                         // 修改时间
	PlanId               int64       `json:"planId"               description:"计划ID"`                                         // 计划ID
	ChannelId            int64       `json:"channelId"            description:"支付渠道Id"`                                       // 支付渠道Id
	Status               int         `json:"status"               description:"渠道绑定状态，0-Init | 1-Create｜2-Active｜3-Inactive"` // 渠道绑定状态，0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelPlanId        string      `json:"channelPlanId"        description:"支付渠道plan_Id"`                                  // 支付渠道plan_Id
	ChannelProductId     string      `json:"channelProductId"     description:"支付渠道product_Id"`                               // 支付渠道product_Id
	ChannelPlanStatus    string      `json:"channelPlanStatus"    description:"channel_plan_status"`                          // channel_plan_status
	ChannelProductStatus string      `json:"channelProductStatus" description:"channel_product_status"`                       // channel_product_status
	Data                 string      `json:"data"                 description:"渠道额外参数，JSON格式"`                                // 渠道额外参数，JSON格式
	IsDeleted            int         `json:"isDeleted"            description:""`                                             //
}
