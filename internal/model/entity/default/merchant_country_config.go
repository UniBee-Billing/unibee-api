// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantCountryConfig is the golang structure for table merchant_country_config.
type MerchantCountryConfig struct {
	Id          int64       `json:"id"          description:""`                      //
	MerchantId  uint64      `json:"merchantId"  description:""`                      //
	CountryCode string      `json:"countryCode" description:""`                      //
	Name        string      `json:"name"        description:""`                      //
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update time"`           // update time
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime  int64       `json:"createTime"  description:"create utc time"`       // create utc time
}
