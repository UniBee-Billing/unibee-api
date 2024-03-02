// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantMember is the golang structure of table merchant_member for DAO operations like Where/Data.
type MerchantMember struct {
	g.Meta     `orm:"table:merchant_member, do:true"`
	Id         interface{} // userId
	GmtCreate  *gtime.Time // create time
	GmtModify  *gtime.Time // update time
	MerchantId interface{} // merchant id
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
	Password   interface{} // password
	UserName   interface{} // user name
	Mobile     interface{} // mobile
	Email      interface{} // email
	FirstName  interface{} // first name
	LastName   interface{} // last name
	CreateTime interface{} // create utc time
}
