// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelVatRate is the golang structure for table channel_vat_rate.
type ChannelVatRate struct {
	Id               uint64      `json:"id"               description:""`                      //
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`           // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`           // update time
	VatRateId        int64       `json:"vatRateId"        description:"vat_rate_id"`           // vat_rate_id
	ChannelId        int64       `json:"channelId"        description:"channel_id"`            // channel_id
	ChannelVatRateId string      `json:"channelVatRateId" description:"channel_vat_rate_Id"`   // channel_vat_rate_Id
	IsDeleted        int         `json:"isDeleted"        description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
}
