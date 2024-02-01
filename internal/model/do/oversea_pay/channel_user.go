// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelUser is the golang structure of table channel_user for DAO operations like Where/Data.
type ChannelUser struct {
	g.Meta                      `orm:"table:channel_user, do:true"`
	Id                          interface{} //
	GmtCreate                   *gtime.Time // create time
	GmtModify                   *gtime.Time // update time
	UserId                      interface{} // userId
	ChannelId                   interface{} // 支付渠道Id
	ChannelUserId               interface{} // 支付渠道user_Id
	IsDeleted                   interface{} // 0-UnDeleted，1-Deleted
	ChannelDefaultPaymentMethod interface{} //
}
