// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelHttpLog is the golang structure for table channel_http_log.
type ChannelHttpLog struct {
	Id        uint64      `json:"id"        description:"id"`               // id
	Url       string      `json:"url"       description:"请求url"`            // 请求url
	Request   string      `json:"request"   description:"请求body参数(json格式)"` // 请求body参数(json格式)
	Response  string      `json:"response"  description:"请求返回结果(json格式)"`   // 请求返回结果(json格式)
	RequestId string      `json:"requestId" description:"reuqest_id"`       // reuqest_id
	Mamo      string      `json:"mamo"      description:"备注"`               // 备注
	ChannelId string      `json:"channelId" description:"channel_id"`       // channel_id
	GmtCreate *gtime.Time `json:"gmtCreate" description:"创建时间"`             // 创建时间
	GmtModify *gtime.Time `json:"gmtModify" description:"更新时间"`             // 更新时间
}
