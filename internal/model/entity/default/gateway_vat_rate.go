// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayVatRate is the golang structure for table gateway_vat_rate.
type GatewayVatRate struct {
	Id               uint64      `json:"id"               description:""`                      //
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`           // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"update time"`           // update time
	VatRateId        int64       `json:"vatRateId"        description:"vat_rate_id"`           // vat_rate_id
	GatewayId        uint64      `json:"gatewayId"        description:"gateway_id"`            // gateway_id
	GatewayVatRateId string      `json:"gatewayVatRateId" description:"gateway_vat_rate_Id"`   // gateway_vat_rate_Id
	IsDeleted        int         `json:"isDeleted"        description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime       int64       `json:"createTime"       description:"create utc time"`       // create utc time
}
