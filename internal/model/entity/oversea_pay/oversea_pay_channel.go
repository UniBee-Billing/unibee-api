// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OverseaPayChannel is the golang structure for table oversea_pay_channel.
type OverseaPayChannel struct {
	Id               uint64      `json:"id"               ` // 主键id
	EnumKey          int64       `json:"enumKey"          ` // 支付渠道枚举（内部定义）
	ChannelType      int         `json:"channelType"      ` // 支付渠道类型，null或者 0-Payment 类型 ｜ 1-Subscription类型
	Channel          string      `json:"channel"          ` // 支付方式枚举（渠道定义）
	Name             string      `json:"name"             ` // 支付方式名称
	SubChannel       string      `json:"subChannel"       ` // 渠道子支付方式枚举
	BrandData        string      `json:"brandData"        ` //
	Logo             string      `json:"logo"             ` // 支付方式logo
	Host             string      `json:"host"             ` // pay host
	ChannelAccountId string      `json:"channelAccountId" ` // 渠道账户Id
	ChannelKey       string      `json:"channelKey"       ` //
	ChannelSecret    string      `json:"channelSecret"    ` // secret
	Custom           string      `json:"custom"           ` // custom
	GmtCreate        *gtime.Time `json:"gmtCreate"        ` // 创建时间
	GmtModify        *gtime.Time `json:"gmtModify"        ` // 修改时间
	Description      string      `json:"description"      ` // 支付方式描述
	WebhookKey       string      `json:"webhookKey"       ` // webhook_key
	WebhookSecret    string      `json:"webhookSecret"    ` // webhook_secret
	UniqueProductId  string      `json:"uniqueProductId"  ` // 渠道唯一 ProductId，目前仅限 Paypal 使用
}
