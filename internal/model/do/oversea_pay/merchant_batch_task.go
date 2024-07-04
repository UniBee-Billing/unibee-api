// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchTask is the golang structure of table merchant_batch_task for DAO operations like Where/Data.
type MerchantBatchTask struct {
	g.Meta         `orm:"table:merchant_batch_task, do:true"`
	Id             interface{} // id
	MerchantId     interface{} // merchant_id
	MemberId       interface{} // member_id
	ModuleName     interface{} // module_name
	TaskName       interface{} // task_name
	SuccessCount   interface{} // success_count
	LastUpdateTime interface{} // last update utc time
	SourceFrom     interface{} // source_from
	Payload        interface{} // payload(json)
	DownloadUrl    interface{} // download_file_url
	Status         interface{} // Status。0-Pending，1-Processing，2-Success，3-Failure
	StartTime      interface{} // task_start_time
	FinishTime     interface{} // task_finish_time
	TaskCost       interface{} // task cost time(second)
	FailReason     interface{} // reason of failure
	GmtCreate      *gtime.Time // gmt_create
	TaskType       interface{} // type，0-download，1-upload
	UploadFileUrl  interface{} // the file url of upload type task
	GmtModify      *gtime.Time // update time
	CreateTime     interface{} // create utc time
}
