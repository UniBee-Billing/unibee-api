// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionVatRateChannel is the golang structure for table subscription_vat_rate_channel.
type SubscriptionVatRateChannel struct {
	Id               uint64      `json:"id"               ` //
	GmtCreate        *gtime.Time `json:"gmtCreate"        ` // 创建时间
	GmtModify        *gtime.Time `json:"gmtModify"        ` // 修改时间
	VatRateId        int64       `json:"vatRateId"        ` // vat_rate_id
	ChannelId        int64       `json:"channelId"        ` // 支付渠道Id
	ChannelVatRateId string      `json:"channelVatRateId" ` // 支付渠道vat_rate_Id
	IsDeleted        int         `json:"isDeleted"        ` //
}
