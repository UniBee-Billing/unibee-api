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
	MemberId           interface{} // member_id
	UserId             interface{} // user_id
	OptAccount         interface{} // admin account
	OptAccountType     interface{} // opt_account_type, 0-Member|1-User|2-OpenApi|3-System
	OptAccountId       interface{} // operation account id
	ClientType         interface{} // client type
	BizType            interface{} // biz_type
	OptTarget          interface{} // operation target
	OptContent         interface{} // operation content
	CreateTime         interface{} // operation create utc time
	IsDelete           interface{} // 0-UnDeleted，1-Deleted
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	QueryportRequestId interface{} // queryport id
	ServerType         interface{} // server type
	ServerTypeDesc     interface{} // server type description
	SubscriptionId     interface{} // subscription_id
	InvoiceId          interface{} // invoice id
	PlanId             interface{} // plan id
	DiscountCode       interface{} // discount_code
}
