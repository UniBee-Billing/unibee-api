// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchTaskExportChunks is the golang structure for table merchant_batch_task_export_chunks.
type MerchantBatchTaskExportChunks struct {
	Id         uint64      `json:"id"         description:"id"`                    // id
	MerchantId uint64      `json:"merchantId" description:"merchant_id"`           // merchant_id
	TaskId     uint64      `json:"taskId"     description:"task_id"`               // task_id
	Page       int         `json:"page"       description:"page"`                  // page
	Content    string      `json:"content"    description:"content"`               // content
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"gmt_create"`            // gmt_create
	GmtModify  *gtime.Time `json:"gmtModify"  description:"update time"`           // update time
	CreateTime int64       `json:"createTime" description:"create utc time"`       // create utc time
}
