// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelUser is the golang structure for table channel_user.
type ChannelUser struct {
	Id                          uint64      `json:"id"                          description:""`                      //
	GmtCreate                   *gtime.Time `json:"gmtCreate"                   description:"创建时间"`                  // 创建时间
	GmtModify                   *gtime.Time `json:"gmtModify"                   description:"修改时间"`                  // 修改时间
	UserId                      int64       `json:"userId"                      description:"userId"`                // userId
	ChannelId                   int64       `json:"channelId"                   description:"支付渠道Id"`                // 支付渠道Id
	ChannelUserId               string      `json:"channelUserId"               description:"支付渠道user_Id"`           // 支付渠道user_Id
	IsDeleted                   int         `json:"isDeleted"                   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	ChannelDefaultPaymentMethod string      `json:"channelDefaultPaymentMethod" description:""`                      //
}