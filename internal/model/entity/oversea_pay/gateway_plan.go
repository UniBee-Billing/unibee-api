// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayPlan is the golang structure for table gateway_plan.
type GatewayPlan struct {
	Id                   uint64      `json:"id"                   description:""`                                      //
	GmtCreate            *gtime.Time `json:"gmtCreate"            description:"create time"`                           // create time
	GmtModify            *gtime.Time `json:"gmtModify"            description:"update time"`                           // update time
	PlanId               int64       `json:"planId"               description:"PlanId"`                                // PlanId
	GatewayId            int64       `json:"gatewayId"            description:"gateway_id"`                            // gateway_id
	Status               int         `json:"status"               description:"0-Init | 1-Create｜2-Active｜3-Inactive"` // 0-Init | 1-Create｜2-Active｜3-Inactive
	GatewayPlanId        string      `json:"gatewayPlanId"        description:"gateway_plan_id"`                       // gateway_plan_id
	GatewayProductId     string      `json:"gatewayProductId"     description:"gateway_product_id"`                    // gateway_product_id
	GatewayPlanStatus    string      `json:"gatewayPlanStatus"    description:"gateway_plan_status"`                   // gateway_plan_status
	GatewayProductStatus string      `json:"gatewayProductStatus" description:"gateway_product_status"`                // gateway_product_status
	IsDeleted            int         `json:"isDeleted"            description:"0-UnDeleted，1-Deleted"`                 // 0-UnDeleted，1-Deleted
	Data                 string      `json:"data"                 description:"data(json)"`                            // data(json)
}
