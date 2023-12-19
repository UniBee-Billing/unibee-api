// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionPlanChannel is the golang structure of table subscription_plan_channel for DAO operations like Where/Data.
type SubscriptionPlanChannel struct {
	g.Meta           `orm:"table:subscription_plan_channel, do:true"`
	Id               interface{} //
	GmtCreate        *gtime.Time // 创建时间
	GmtModify        *gtime.Time // 修改时间
	PlanId           interface{} // 计划ID
	ChannelId        interface{} // 支付渠道Id
	ChannelPlanId    interface{} // 支付渠道plan_Id
	ChannelProductId interface{} // 支付渠道product_Id
	Data             interface{} // 渠道额外参数，JSON格式
	IsDeleted        interface{} //
}
