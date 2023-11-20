// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelHttpLog is the golang structure for table channel_http_log.
type ChannelHttpLog struct {
	Id        uint64      `json:"id"        ` // id
	Url       string      `json:"url"       ` // 请求url
	Request   string      `json:"request"   ` // 请求body参数(json格式)
	Response  string      `json:"response"  ` // 请求返回结果(json格式)
	RequestId string      `json:"requestId" ` // reuqest_id
	Mamo      string      `json:"mamo"      ` // 备注
	ChannelId string      `json:"channelId" ` // channel_id
	GmtCreate *gtime.Time `json:"gmtCreate" ` // 创建时间
	GmtModify *gtime.Time `json:"gmtModify" ` // 更新时间
}
