// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelHttpLog is the golang structure for table channel_http_log.
type ChannelHttpLog struct {
	Id        uint64      `json:"id"        description:"id"`                 // id
	Url       string      `json:"url"       description:"request url"`        // request url
	Request   string      `json:"request"   description:"request body(json)"` // request body(json)
	Response  string      `json:"response"  description:"response(json)"`     // response(json)
	RequestId string      `json:"requestId" description:"reuqest_id"`         // reuqest_id
	Mamo      string      `json:"mamo"      description:"mamo"`               // mamo
	ChannelId string      `json:"channelId" description:"channel_id"`         // channel_id
	GmtCreate *gtime.Time `json:"gmtCreate" description:"create time"`        // create time
	GmtModify *gtime.Time `json:"gmtModify" description:"update time"`        // update time
}
