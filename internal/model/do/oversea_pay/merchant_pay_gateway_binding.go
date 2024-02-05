// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantPayGatewayBinding is the golang structure of table merchant_pay_gateway_binding for DAO operations like Where/Data.
type MerchantPayGatewayBinding struct {
	g.Meta     `orm:"table:merchant_pay_gateway_binding, do:true"`
	Id         interface{} //
	GmtCreate  *gtime.Time // create time
	GmtModify  *gtime.Time // update time
	MerchantId interface{} // merchant id
	GatewayId  interface{} // gateway_id
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
}
