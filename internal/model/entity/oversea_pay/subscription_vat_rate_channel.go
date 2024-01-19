// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionVatRateChannel is the golang structure for table subscription_vat_rate_channel.
type SubscriptionVatRateChannel struct {
	Id               uint64      `json:"id"               description:""`                //
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"创建时间"`            // 创建时间
	GmtModify        *gtime.Time `json:"gmtModify"        description:"修改时间"`            // 修改时间
	VatRateId        int64       `json:"vatRateId"        description:"vat_rate_id"`     // vat_rate_id
	ChannelId        int64       `json:"channelId"        description:"支付渠道Id"`          // 支付渠道Id
	ChannelVatRateId string      `json:"channelVatRateId" description:"支付渠道vat_rate_Id"` // 支付渠道vat_rate_Id
	IsDeleted        int         `json:"isDeleted"        description:""`                //
}
