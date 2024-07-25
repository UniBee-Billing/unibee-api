// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRole is the golang structure of table merchant_role for DAO operations like Where/Data.
type MerchantRole struct {
	g.Meta         `orm:"table:merchant_role, do:true"`
	Id             interface{} // userId
	GmtCreate      *gtime.Time // create time
	GmtModify      *gtime.Time // update time
	MerchantId     interface{} // merchant id
	IsDeleted      interface{} // 0-UnDeleted，1-Deleted
	Role           interface{} // role
	PermissionData interface{} // permission_data（json）
	CreateTime     interface{} // create utc time
}
