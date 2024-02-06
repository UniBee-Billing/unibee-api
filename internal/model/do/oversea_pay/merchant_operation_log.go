// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantOperationLog is the golang structure of table merchant_operation_log for DAO operations like Where/Data.
type MerchantOperationLog struct {
	g.Meta             `orm:"table:merchant_operation_log, do:true"`
	Id                 interface{} // id
	CompanyId          interface{} // company id
	MerchantId         interface{} // merchant Id
	UserId             interface{} // user_id
	OptAccount         interface{} // admin account
	ClientType         interface{} // client type
	BizType            interface{} // biz_type
	OptTarget          interface{} // operation target
	OptContent         interface{} // operation content
	CreateAt           interface{} // operation create time
	IsDelete           interface{} // 0-UnDeleted，1-Deleted
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	QueryportRequestId interface{} // queryport id
	ServerType         interface{} // server type
	ServerTypeDesc     interface{} // server type description
}
