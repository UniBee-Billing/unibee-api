// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantRole is the golang structure for table merchant_role.
type MerchantRole struct {
	Id             uint64      `json:"id"             description:"userId"`                // userId
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"create time"`           // create time
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`           // update time
	MerchantId     uint64      `json:"merchantId"     description:"merchant id"`           // merchant id
	IsDeleted      int         `json:"isDeleted"      description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	Role           string      `json:"role"           description:"role"`                  // role
	PermissionData string      `json:"permissionData" description:"permission_data（json）"` // permission_data（json）
	CreateTime     int64       `json:"createTime"     description:"create utc time"`       // create utc time
}
