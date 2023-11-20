// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantPayChannelMapping is the golang structure for table merchant_pay_channel_mapping.
type MerchantPayChannelMapping struct {
	Id         uint64      `json:"id"         ` //
	GmtCreate  *gtime.Time `json:"gmtCreate"  ` // 创建时间
	GmtModify  *gtime.Time `json:"gmtModify"  ` // 修改时间
	MerchantId int64       `json:"merchantId" ` // 商户Id
	ChannelId  string      `json:"channelId"  ` // oversea_pay_channel表的id
	IsDeleted  int         `json:"isDeleted"  ` //
}
