// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAdminNote is the golang structure of table user_admin_note for DAO operations like Where/Data.
type UserAdminNote struct {
	g.Meta           `orm:"table:user_admin_note, do:true"`
	Id               interface{} // id
	GmtCreate        *gtime.Time // create_time
	GmtModify        *gtime.Time // modify_time
	UserId           interface{} // user_id
	MerchantMemberId interface{} // merchant_user_id
	Note             interface{} // note
	IsDeleted        interface{} // 0-UnDeleted，1-Deleted
	CreateTime       interface{} // create utc time
}
