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
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"创建时间"`                    // 创建时间
	GmtModify  *gtime.Time `json:"gmtModify"  description:"修改时间"`                    // 修改时间
	MerchantId int64       `json:"merchantId" description:"商户Id"`                    // 商户Id
	ChannelId  string      `json:"channelId"  description:"oversea_pay_channel表的id"` // oversea_pay_channel表的id
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"`   // 0-UnDeleted，1-Deleted
}
