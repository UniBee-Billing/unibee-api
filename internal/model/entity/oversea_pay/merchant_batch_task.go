// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantBatchTask is the golang structure for table merchant_batch_task.
type MerchantBatchTask struct {
	Id             int64       `json:"id"             description:"id"`                                                // id
	MerchantId     uint64      `json:"merchantId"     description:"merchant_id"`                                       // merchant_id
	MemberId       uint64      `json:"memberId"       description:"member_id"`                                         // member_id
	ModuleName     string      `json:"moduleName"     description:"module_name"`                                       // module_name
	TaskName       string      `json:"taskName"       description:"task_name"`                                         // task_name
	SuccessCount   int64       `json:"successCount"   description:"success_count"`                                     // success_count
	LastUpdateTime int64       `json:"lastUpdateTime" description:"last update utc time"`                              // last update utc time
	SourceFrom     string      `json:"sourceFrom"     description:"source_from"`                                       // source_from
	Payload        string      `json:"payload"        description:"payload(json)"`                                     // payload(json)
	DownloadUrl    string      `json:"downloadUrl"    description:"download_file_url"`                                 // download_file_url
	Status         int         `json:"status"         description:"Status。0-Pending，1-Processing，2-Success，3-Failure"` // Status。0-Pending，1-Processing，2-Success，3-Failure
	StartTime      int64       `json:"startTime"      description:"task_start_time"`                                   // task_start_time
	FinishTime     int64       `json:"finishTime"     description:"task_finish_time"`                                  // task_finish_time
	TaskCost       int         `json:"taskCost"       description:"task cost time(second)"`                            // task cost time(second)
	FailReason     string      `json:"failReason"     description:"reason of failure"`                                 // reason of failure
	GmtCreate      *gtime.Time `json:"gmtCreate"      description:"gmt_create"`                                        // gmt_create
	TaskType       int         `json:"taskType"       description:"type，0-download，1-upload"`                          // type，0-download，1-upload
	UploadFileUrl  string      `json:"uploadFileUrl"  description:"the file url of upload type task"`                  // the file url of upload type task
	GmtModify      *gtime.Time `json:"gmtModify"      description:"update time"`                                       // update time
	CreateTime     int64       `json:"createTime"     description:"create utc time"`                                   // create utc time
}
