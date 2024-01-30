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
	GmtCreate  *gtime.Time // 创建时间
	GmtModify  *gtime.Time // 修改时间
	MerchantId interface{} // 商户Id
	ChannelId  interface{} // oversea_pay_channel表的id
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
}
