// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayPlan is the golang structure of table gateway_plan for DAO operations like Where/Data.
type GatewayPlan struct {
	g.Meta               `orm:"table:gateway_plan, do:true"`
	Id                   interface{} //
	GmtCreate            *gtime.Time // create time
	GmtModify            *gtime.Time // update time
	PlanId               interface{} // PlanId
	GatewayId            interface{} // gateway_id
	Status               interface{} // 0-Init | 1-Create｜2-Active｜3-Inactive
	GatewayPlanId        interface{} // gateway_plan_id
	GatewayProductId     interface{} // gateway_product_id
	GatewayPlanStatus    interface{} // gateway_plan_status
	GatewayProductStatus interface{} // gateway_product_status
	IsDeleted            interface{} // 0-UnDeleted，1-Deleted
	Data                 interface{} // data(json)
	CreateAt             interface{} // create utc time
}
