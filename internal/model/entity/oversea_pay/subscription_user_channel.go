// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionUserChannel is the golang structure for table subscription_user_channel.
type SubscriptionUserChannel struct {
	Id            uint64      `json:"id"            ` //
	GmtCreate     *gtime.Time `json:"gmtCreate"     ` // 创建时间
	GmtModify     *gtime.Time `json:"gmtModify"     ` // 修改时间
	UserId        int64       `json:"userId"        ` // userId
	ChannelId     int64       `json:"channelId"     ` // 支付渠道Id
	ChannelUserId string      `json:"channelUserId" ` // 支付渠道user_Id
	IsDeleted     int         `json:"isDeleted"     ` //
}
