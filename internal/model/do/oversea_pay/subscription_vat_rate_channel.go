// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SubscriptionVatRateChannel is the golang structure of table subscription_vat_rate_channel for DAO operations like Where/Data.
type SubscriptionVatRateChannel struct {
	g.Meta           `orm:"table:subscription_vat_rate_channel, do:true"`
	Id               interface{} //
	GmtCreate        *gtime.Time // 创建时间
	GmtModify        *gtime.Time // 修改时间
	VatRateId        interface{} // vat_rate_id
	ChannelId        interface{} // 支付渠道Id
	ChannelVatRateId interface{} // 支付渠道vat_rate_Id
	IsDeleted        interface{} // 0-UnDeleted，1-Deleted
}
