// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionUserChannel is the golang structure of table subscription_user_channel for DAO operations like Where/Data.
type SubscriptionUserChannel struct {
	g.Meta        `orm:"table:subscription_user_channel, do:true"`
	Id            interface{} //
	GmtCreate     *gtime.Time // 创建时间
	GmtModify     *gtime.Time // 修改时间
	UserId        interface{} // userId
	ChannelId     interface{} // 支付渠道Id
	ChannelUserId interface{} // 支付渠道user_Id
	IsDeleted     interface{} //
}
