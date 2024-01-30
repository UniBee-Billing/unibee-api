// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantChannelConfig is the golang structure of table merchant_channel_config for DAO operations like Where/Data.
type MerchantChannelConfig struct {
	g.Meta           `orm:"table:merchant_channel_config, do:true"`
	Id               interface{} // 主键id
	MerchantId       interface{} //
	EnumKey          interface{} // 支付渠道枚举（内部定义）
	ChannelType      interface{} // 支付渠道类型，null或者 0-Payment 类型 ｜ 1-Subscription类型
	Channel          interface{} // 支付方式枚举（渠道定义）
	Name             interface{} // 支付方式名称
	SubChannel       interface{} // 渠道子支付方式枚举
	BrandData        interface{} //
	Logo             interface{} // 支付方式logo
	Host             interface{} // pay host
	ChannelAccountId interface{} // 渠道账户Id
	ChannelKey       interface{} //
	ChannelSecret    interface{} // secret
	Custom           interface{} // custom
	GmtCreate        *gtime.Time // 创建时间
	GmtModify        *gtime.Time // 修改时间
	Description      interface{} // 支付方式描述
	WebhookKey       interface{} // webhook_key
	WebhookSecret    interface{} // webhook_secret
	UniqueProductId  interface{} // 渠道唯一 ProductId，目前仅限 Paypal 使用
}
