// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OpenApiConfig is the golang structure for table open_api_config.
type OpenApiConfig struct {
	Id              uint64      `json:"id"              ` //
	Qps             int         `json:"qps"             ` // 开放平台Api qps总控制
	GmtCreate       *gtime.Time `json:"gmtCreate"       ` // 创建时间
	GmtModify       *gtime.Time `json:"gmtModify"       ` // 修改时间
	MerchantId      int64       `json:"merchantId"      ` // 商户Id
	Hmac            string      `json:"hmac"            ` // 回调加密hmac
	Callback        string      `json:"callback"        ` // 回调Url
	ApiKey          string      `json:"apiKey"          ` // 开放平台Key
	Token           string      `json:"token"           ` // 开放平台token
	IsDeleted       int         `json:"isDeleted"       ` //
	Validips        string      `json:"validips"        ` //
	ChannelCallback string      `json:"channelCallback" ` // 渠道支付原信息回调地址
	CompanyId       int64       `json:"companyId"       ` // 公司ID
}
