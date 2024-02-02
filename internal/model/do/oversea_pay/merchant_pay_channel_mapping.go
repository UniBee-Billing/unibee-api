// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantPayChannelMapping is the golang structure of table merchant_pay_channel_mapping for DAO operations like Where/Data.
type MerchantPayChannelMapping struct {
	g.Meta     `orm:"table:merchant_pay_channel_mapping, do:true"`
	Id         interface{} //
	GmtCreate  *gtime.Time // create time
	GmtModify  *gtime.Time // update time
	MerchantId interface{} // merchant id
	ChannelId  interface{} // oversea_pay_channel表的id
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
}
