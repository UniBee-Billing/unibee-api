// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelPlan is the golang structure of table channel_plan for DAO operations like Where/Data.
type ChannelPlan struct {
	g.Meta               `orm:"table:channel_plan, do:true"`
	Id                   interface{} //
	GmtCreate            *gtime.Time // create time
	GmtModify            *gtime.Time // update time
	PlanId               interface{} // PlanId
	ChannelId            interface{} // 支付渠道Id
	Status               interface{} // 0-Init | 1-Create｜2-Active｜3-Inactive
	ChannelPlanId        interface{} // 支付渠道plan_Id
	ChannelProductId     interface{} // 支付渠道product_Id
	ChannelPlanStatus    interface{} // channel_plan_status
	ChannelProductStatus interface{} // channel_product_status
	IsDeleted            interface{} // 0-UnDeleted，1-Deleted
	Data                 interface{} // 渠道额外参数，JSON格式
}
