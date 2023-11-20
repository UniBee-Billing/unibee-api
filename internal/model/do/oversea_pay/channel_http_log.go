// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChannelHttpLog is the golang structure of table channel_http_log for DAO operations like Where/Data.
type ChannelHttpLog struct {
	g.Meta    `orm:"table:channel_http_log, do:true"`
	Id        interface{} // id
	Url       interface{} // 请求url
	Request   interface{} // 请求body参数(json格式)
	Response  interface{} // 请求返回结果(json格式)
	RequestId interface{} // reuqest_id
	Mamo      interface{} // 备注
	ChannelId interface{} // channel_id
	GmtCreate *gtime.Time // 创建时间
	GmtModify *gtime.Time // 更新时间
}
