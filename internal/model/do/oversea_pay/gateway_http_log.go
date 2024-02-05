// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GatewayHttpLog is the golang structure of table gateway_http_log for DAO operations like Where/Data.
type GatewayHttpLog struct {
	g.Meta    `orm:"table:gateway_http_log, do:true"`
	Id        interface{} // id
	Url       interface{} // request url
	Request   interface{} // request body(json)
	Response  interface{} // response(json)
	RequestId interface{} // request_id
	Mamo      interface{} // mamo
	GatewayId interface{} // gateway_id
	GmtCreate *gtime.Time // create time
	GmtModify *gtime.Time // update time
}
