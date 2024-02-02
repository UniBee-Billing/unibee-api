// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantConfig is the golang structure for table merchant_config.
type MerchantConfig struct {
	Id          int64       `json:"id"          description:"ID"`                    // ID
	MerchantId  int64       `json:"merchantId"  description:"merchantId"`            // merchantId
	ConfigKey   string      `json:"configKey"   description:"config_key"`            // config_key
	ConfigValue string      `json:"configValue" description:"config_value"`          // config_value
	GmtCreate   *gtime.Time `json:"gmtCreate"   description:"create time"`           // create time
	GmtModify   *gtime.Time `json:"gmtModify"   description:"update time"`           // update time
	IsDeleted   int         `json:"isDeleted"   description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
}
