// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantEmailHistory is the golang structure of table merchant_email_history for DAO operations like Where/Data.
type MerchantEmailHistory struct {
	g.Meta     `orm:"table:merchant_email_history, do:true"`
	Id         interface{} //
	MerchantId interface{} //
	Email      interface{} //
	Title      interface{} //
	Content    interface{} //
	AttachFile interface{} //
	GmtCreate  *gtime.Time // create time
	GmtModify  *gtime.Time // update time
	Response   interface{} //
	CreateTime interface{} // create utc time
}
