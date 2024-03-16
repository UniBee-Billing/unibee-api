// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantCountryConfig is the golang structure of table merchant_country_config for DAO operations like Where/Data.
type MerchantCountryConfig struct {
	g.Meta      `orm:"table:merchant_country_config, do:true"`
	Id          interface{} //
	MerchantId  interface{} //
	CountryCode interface{} //
	Name        interface{} //
	GmtCreate   *gtime.Time // create time
	GmtModify   *gtime.Time // update time
	IsDeleted   interface{} // 0-UnDeleted，1-Deleted
	CreateTime  interface{} // create utc time
}
