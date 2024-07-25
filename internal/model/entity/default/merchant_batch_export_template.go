// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchExportTemplate is the golang structure for table merchant_batch_export_template.
type MerchantBatchExportTemplate struct {
	Id            uint64      `json:"id"            description:"id"`                    // id
	MerchantId    uint64      `json:"merchantId"    description:"merchant_id"`           // merchant_id
	MemberId      uint64      `json:"memberId"      description:"member_id"`             // member_id
	Name          string      `json:"name"          description:"name"`                  // name
	Task          string      `json:"task"          description:"task"`                  // task
	Format        string      `json:"format"        description:"format"`                // format
	Payload       string      `json:"payload"       description:"payload(json)"`         // payload(json)
	ExportColumns string      `json:"exportColumns" description:"export_columns(json)"`  // export_columns(json)
	IsDeleted     int         `json:"isDeleted"     description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	GmtCreate     *gtime.Time `json:"gmtCreate"     description:"gmt_create"`            // gmt_create
	GmtModify     *gtime.Time `json:"gmtModify"     description:"update time"`           // update time
	CreateTime    int64       `json:"createTime"    description:"create utc time"`       // create utc time
}
