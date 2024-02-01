// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OpenApiConfig is the golang structure for table open_api_config.
type OpenApiConfig struct {
	Id              uint64      `json:"id"              description:""`                      //
	Qps             int         `json:"qps"             description:"开放平台Api qps总控制"`        // 开放平台Api qps总控制
	GmtCreate       *gtime.Time `json:"gmtCreate"       description:"create time"`           // create time
	GmtModify       *gtime.Time `json:"gmtModify"       description:"修改时间"`                  // 修改时间
	MerchantId      int64       `json:"merchantId"      description:"商户Id"`                  // 商户Id
	Hmac            string      `json:"hmac"            description:"回调加密hmac"`              // 回调加密hmac
	Callback        string      `json:"callback"        description:"回调Url"`                 // 回调Url
	ApiKey          string      `json:"apiKey"          description:"开放平台Key"`               // 开放平台Key
	Token           string      `json:"token"           description:"开放平台token"`             // 开放平台token
	IsDeleted       int         `json:"isDeleted"       description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Validips        string      `json:"validips"        description:""`                      //
	ChannelCallback string      `json:"channelCallback" description:"渠道支付原信息回调地址"`           // 渠道支付原信息回调地址
	CompanyId       int64       `json:"companyId"       description:"公司ID"`                  // 公司ID
}
