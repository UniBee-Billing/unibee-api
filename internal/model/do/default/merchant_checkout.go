// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantCheckout is the golang structure of table merchant_checkout for DAO operations like Where/Data.
type MerchantCheckout struct {
	g.Meta      `orm:"table:merchant_checkout, do:true"`
	Id          interface{} // ID
	MerchantId  interface{} // merchantId
	Name        interface{} // name
	Description interface{} // description
	Data        interface{} // data(json)
	Staging     interface{} // staging_data(json)
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update time
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CreateTime  interface{} // create utc time
}
