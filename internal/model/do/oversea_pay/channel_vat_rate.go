// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelVatRate is the golang structure of table channel_vat_rate for DAO operations like Where/Data.
type ChannelVatRate struct {
	g.Meta           `orm:"table:channel_vat_rate, do:true"`
	Id               interface{} //
	GmtCreate        *gtime.Time // create time
	GmtModify        *gtime.Time // update time
	VatRateId        interface{} // vat_rate_id
	ChannelId        interface{} // 支付渠道Id
	ChannelVatRateId interface{} // 支付渠道vat_rate_Id
	IsDeleted        interface{} // 0-UnDeleted，1-Deleted
}
