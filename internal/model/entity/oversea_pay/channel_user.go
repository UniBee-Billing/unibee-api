// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelUser is the golang structure for table channel_user.
type ChannelUser struct {
	Id                          uint64      `json:"id"                          description:""`                               //
	GmtCreate                   *gtime.Time `json:"gmtCreate"                   description:"create time"`                    // create time
	GmtModify                   *gtime.Time `json:"gmtModify"                   description:"update time"`                    // update time
	UserId                      int64       `json:"userId"                      description:"userId"`                         // userId
	ChannelId                   int64       `json:"channelId"                   description:"channel_id"`                     // channel_id
	ChannelUserId               string      `json:"channelUserId"               description:"channel_user_Id"`                // channel_user_Id
	IsDeleted                   int         `json:"isDeleted"                   description:"0-UnDeleted，1-Deleted"`          // 0-UnDeleted，1-Deleted
	ChannelDefaultPaymentMethod string      `json:"channelDefaultPaymentMethod" description:"channel_default_payment_method"` // channel_default_payment_method
}
