// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionAdminNote is the golang structure of table subscription_admin_note for DAO operations like Where/Data.
type SubscriptionAdminNote struct {
	g.Meta         `orm:"table:subscription_admin_note, do:true"`
	Id             interface{} // id
	GmtCreate      *gtime.Time // create_time
	GmtModify      *gtime.Time // modify_time
	SubscriptionId interface{} // subscription_id
	MerchantUserId interface{} // merchant_user_id
	Note           interface{} // note
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
}
