// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPayChannel is the golang structure of table oversea_pay_channel for DAO operations like Where/Data.
type OverseaPayChannel struct {
	g.Meta           `orm:"table:oversea_pay_channel, do:true"`
	Id               interface{} // 主键id
	EnumKey          interface{} // 支付渠道枚举（内部定义）
	Channel          interface{} // 支付方式枚举（渠道定义）
	Name             interface{} // 支付方式名称
	SubChannel       interface{} // 渠道子支付方式枚举
	BrandData        interface{} //
	Logo             interface{} // 支付方式logo
	ChannelAccountId interface{} // 渠道账户Id
	ChannelKey       interface{} //
	ChannelSecret    interface{} // secret
	Custom           interface{} // custom
	GmtCreate        *gtime.Time // 创建时间
	GmtModify        *gtime.Time // 修改时间
	Description      interface{} // 支付方式描述
}
