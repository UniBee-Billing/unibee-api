// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantConfig is the golang structure for table merchant_config.
type MerchantConfig struct {
	Id          int64       `json:"id"          ` // ID
	MerchantId  int64       `json:"merchantId"  ` // merchantId
	ConfigKey   string      `json:"configKey"   ` // config_key
	ConfigValue string      `json:"configValue" ` // config_value
	GmtCreate   *gtime.Time `json:"gmtCreate"   ` // 创建时间
	GmtModify   *gtime.Time `json:"gmtModify"   ` // 修改时间
	IsDeleted   int         `json:"isDeleted"   ` // 是否删除，0-未删除，1-删除
}
