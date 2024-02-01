// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantChannelConfig is the golang structure for table merchant_channel_config.
type MerchantChannelConfig struct {
	Id               uint64      `json:"id"               description:"主键id"`                                          // 主键id
	MerchantId       int64       `json:"merchantId"       description:""`                                              //
	EnumKey          int64       `json:"enumKey"          description:"支付渠道枚举（内部定义）"`                                  // 支付渠道枚举（内部定义）
	ChannelType      int         `json:"channelType"      description:"支付渠道类型，null或者 0-Payment 类型 ｜ 1-Subscription类型"` // 支付渠道类型，null或者 0-Payment 类型 ｜ 1-Subscription类型
	Channel          string      `json:"channel"          description:"支付方式枚举（渠道定义）"`                                  // 支付方式枚举（渠道定义）
	Name             string      `json:"name"             description:"支付方式名称"`                                        // 支付方式名称
	SubChannel       string      `json:"subChannel"       description:"渠道子支付方式枚举"`                                     // 渠道子支付方式枚举
	BrandData        string      `json:"brandData"        description:""`                                              //
	Logo             string      `json:"logo"             description:"支付方式logo"`                                      // 支付方式logo
	Host             string      `json:"host"             description:"pay host"`                                      // pay host
	ChannelAccountId string      `json:"channelAccountId" description:"渠道账户Id"`                                        // 渠道账户Id
	ChannelKey       string      `json:"channelKey"       description:""`                                              //
	ChannelSecret    string      `json:"channelSecret"    description:"secret"`                                        // secret
	Custom           string      `json:"custom"           description:"custom"`                                        // custom
	GmtCreate        *gtime.Time `json:"gmtCreate"        description:"create time"`                                   // create time
	GmtModify        *gtime.Time `json:"gmtModify"        description:"修改时间"`                                          // 修改时间
	Description      string      `json:"description"      description:"支付方式描述"`                                        // 支付方式描述
	WebhookKey       string      `json:"webhookKey"       description:"webhook_key"`                                   // webhook_key
	WebhookSecret    string      `json:"webhookSecret"    description:"webhook_secret"`                                // webhook_secret
	UniqueProductId  string      `json:"uniqueProductId"  description:"渠道唯一 ProductId，目前仅限 Paypal 使用"`                 // 渠道唯一 ProductId，目前仅限 Paypal 使用
}
