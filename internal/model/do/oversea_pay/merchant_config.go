// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantConfig is the golang structure of table merchant_config for DAO operations like Where/Data.
type MerchantConfig struct {
	g.Meta      `orm:"table:merchant_config, do:true"`
	Id          interface{} // ID
	MerchantId  interface{} // merchantId
	ConfigKey   interface{} // config_key
	ConfigValue interface{} // config_value
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update time
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
}
