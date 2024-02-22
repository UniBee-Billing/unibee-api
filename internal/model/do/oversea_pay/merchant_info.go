// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantInfo is the golang structure of table merchant_info for DAO operations like Where/Data.
type MerchantInfo struct {
	g.Meta      `orm:"table:merchant_info, do:true"`
	Id          interface{} // merchant_id
	CompanyId   interface{} // company_id
	UserId      interface{} // create_user_id
	Type        interface{} // type
	CompanyName interface{} // company_name
	Email       interface{} // email
	BusinessNum interface{} // business_num
	Name        interface{} // name
	Idcard      interface{} // idcard
	Location    interface{} // location
	Address     interface{} // address
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update_time
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CompanyLogo interface{} // company_logo
	HomeUrl     interface{} //
	Phone       interface{} // phone
	CreateTime  interface{} // create utc time
	TimeZone    interface{} // merchant default time zone
	Host        interface{} // merchant user portal host
	ApiKey      interface{} // merchant open api key
}
