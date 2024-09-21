// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayHttpLog is the golang structure for table gateway_http_log.
type GatewayHttpLog struct {
	Id         uint64      `json:"id"         description:"id"`                 // id
	Url        string      `json:"url"        description:"request url"`        // request url
	Request    string      `json:"request"    description:"request body(json)"` // request body(json)
	Response   string      `json:"response"   description:"response(json)"`     // response(json)
	RequestId  string      `json:"requestId"  description:"request_id"`         // request_id
	Mamo       string      `json:"mamo"       description:"mamo"`               // mamo
	GatewayId  string      `json:"gatewayId"  description:"gateway_id"`         // gateway_id
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"`        // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"`        // update time
	CreateTime int64       `json:"createTime" description:"create utc time"`    // create utc time
}
