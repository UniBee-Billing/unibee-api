// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayUser is the golang structure of table gateway_user for DAO operations like Where/Data.
type GatewayUser struct {
	g.Meta                      `orm:"table:gateway_user, do:true"`
	Id                          interface{} //
	GmtCreate                   *gtime.Time // create time
	GmtModify                   *gtime.Time // update time
	UserId                      interface{} // userId
	GatewayId                   interface{} // gateway_id
	GatewayUserId               interface{} // gateway_user_Id
	IsDeleted                   interface{} // 0-UnDeleted，1-Deleted
	GatewayDefaultPaymentMethod interface{} // gateway_default_payment_method
	CreateAt                    interface{} // create utc time
}
