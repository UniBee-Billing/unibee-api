// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantPayChannelMapping is the golang structure for table merchant_pay_channel_mapping.
type MerchantPayChannelMapping struct {
	Id         uint64      `json:"id"         description:""`                        //
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"`             // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"`             // update time
	MerchantId int64       `json:"merchantId" description:"merchant id"`             // merchant id
	ChannelId  string      `json:"channelId"  description:"oversea_pay_channel表的id"` // oversea_pay_channel表的id
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"`   // 0-UnDeleted，1-Deleted
}
