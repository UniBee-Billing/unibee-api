// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelPlan is the golang structure for table channel_plan.
type ChannelPlan struct {
	Id                   uint64      `json:"id"                   description:""`                                      //
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"创建时间"`                                  // 创建时间
	GmtModify            *gtime.Time `json:"gmtModify"            description:"修改时间"`                                  // 修改时间
	PlanId               int64       `json:"planId"               description:"PlanId"`                                // PlanId
	ChannelId            int64       `json:"channelId"            description:"支付渠道Id"`                                // 支付渠道Id
	Status               int         `json:"status"               description:"0-Init | 1-Create｜2-Active｜3-Inactive"` // 0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelPlanId        string      `json:"channelPlanId"        description:"支付渠道plan_Id"`                           // 支付渠道plan_Id
	ChannelProductId     string      `json:"channelProductId"     description:"支付渠道product_Id"`                        // 支付渠道product_Id
	ChannelPlanStatus    string      `json:"channelPlanStatus"    description:"channel_plan_status"`                   // channel_plan_status
	ChannelProductStatus string      `json:"channelProductStatus" description:"channel_product_status"`                // channel_product_status
	IsDeleted            int         `json:"isDeleted"            description:"0-UnDeleted，1-Deleted"`                 // 0-UnDeleted，1-Deleted
	Data                 string      `json:"data"                 description:"渠道额外参数，JSON格式"`                         // 渠道额外参数，JSON格式
}
