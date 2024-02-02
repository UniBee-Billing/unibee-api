// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OpenApiConfig is the golang structure of table open_api_config for DAO operations like Where/Data.
type OpenApiConfig struct {
	g.Meta          `orm:"table:open_api_config, do:true"`
	Id              interface{} //
	Qps             interface{} // total qps control
	GmtCreate       *gtime.Time // create time
	GmtModify       *gtime.Time // update time
	MerchantId      interface{} // merchant id
	Hmac            interface{} // webhook hmac key
	Callback        interface{} // callback url
	ApiKey          interface{} // api key
	Token           interface{} // api token
	IsDeleted       interface{} // 0-UnDeleted，1-Deleted
	Validips        interface{} //
	ChannelCallback interface{} // callback return data
	CompanyId       interface{} // company id
}
