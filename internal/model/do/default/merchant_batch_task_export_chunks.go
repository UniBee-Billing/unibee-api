// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchTaskExportChunks is the golang structure of table merchant_batch_task_export_chunks for DAO operations like Where/Data.
type MerchantBatchTaskExportChunks struct {
	g.Meta     `orm:"table:merchant_batch_task_export_chunks, do:true"`
	Id         interface{} // id
	MerchantId interface{} // merchant_id
	TaskId     interface{} // task_id
	Page       interface{} // page
	Content    interface{} // content
	IsDeleted  interface{} // 0-UnDeleted，1-Deleted
	GmtCreate  *gtime.Time // gmt_create
	GmtModify  *gtime.Time // update time
	CreateTime interface{} // create utc time
}
