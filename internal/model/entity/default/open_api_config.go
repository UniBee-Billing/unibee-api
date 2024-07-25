// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OpenApiConfig is the golang structure for table open_api_config.
type OpenApiConfig struct {
	Id                      uint64      `json:"id"                      description:""`                         //
	Qps                     int         `json:"qps"                     description:"total qps control"`        // total qps control
	GmtCreate               *gtime.Time `json:"gmtCreate"               description:"create time"`              // create time
	GmtModify               *gtime.Time `json:"gmtModify"               description:"update time"`              // update time
	MerchantId              uint64      `json:"merchantId"              description:"merchant id"`              // merchant id
	Hmac                    string      `json:"hmac"                    description:"webhook hmac key"`         // webhook hmac key
	Callback                string      `json:"callback"                description:"callback url"`             // callback url
	ApiKey                  string      `json:"apiKey"                  description:"api key"`                  // api key
	Token                   string      `json:"token"                   description:"api token"`                // api token
	IsDeleted               int         `json:"isDeleted"               description:"0-UnDeleted，1-Deleted"`    // 0-UnDeleted，1-Deleted
	Validips                string      `json:"validips"                description:""`                         //
	GatewayCallbackResponse string      `json:"gatewayCallbackResponse" description:"callback return response"` // callback return response
	CompanyId               int64       `json:"companyId"               description:"company id"`               // company id
	CreateTime              int64       `json:"createTime"              description:"create utc time"`          // create utc time
}
