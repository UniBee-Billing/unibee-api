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
	Qps             interface{} // 开放平台Api qps总控制
	GmtCreate       *gtime.Time // 创建时间
	GmtModify       *gtime.Time // 修改时间
	MerchantId      interface{} // 商户Id
	Hmac            interface{} // 回调加密hmac
	Callback        interface{} // 回调Url
	ApiKey          interface{} // 开放平台Key
	Token           interface{} // 开放平台token
	IsDeleted       interface{} // 0-UnDeleted，1-Deleted
	Validips        interface{} //
	ChannelCallback interface{} // 渠道支付原信息回调地址
	CompanyId       interface{} // 公司ID
}
