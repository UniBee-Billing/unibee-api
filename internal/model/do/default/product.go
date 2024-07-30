// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Product is the golang structure of table product for DAO operations like Where/Data.
type Product struct {
	g.Meta      `orm:"table:product, do:true"`
	Id          interface{} //
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update time
	CompanyId   interface{} // company id
	MerchantId  interface{} // merchant id
	ProductName interface{} // ProductName
	Description interface{} // description
	ImageUrl    interface{} // image_url
	HomeUrl     interface{} // home_url
	Status      interface{} // status，1-active，2-inactive, default active
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CreateTime  interface{} // create utc time
	MetaData    interface{} // meta_data(json)
}
