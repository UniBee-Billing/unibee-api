package bean

import (
	"net/url"
	"path/filepath"
	"strings"
	"unibee/internal/controller/link"
	entity "unibee/internal/model/entity/oversea_pay"
)

type MerchantBatchTaskSimplify struct {
	Id             int64  `json:"id"            description:"id"`                                                // id
	MerchantId     uint64 `json:"merchantId"    description:"merchant_id"`                                       // merchant_id
	MemberId       uint64 `json:"memberId"      description:"member_id"`                                         // member_id
	TaskName       string `json:"taskName"      description:"task_name"`                                         // task_name
	Payload        string `json:"payload"       description:"payload(json)"`                                     // payload(json)
	DownloadUrl    string `json:"downloadUrl"   description:"download_file_url"`                                 // download_file_url
	Status         int    `json:"status"        description:"Status。0-Pending，1-Processing，2-Success，3-Failure"` // Status。0-Pending，1-Processing，2-Success，3-Failure
	StartTime      int64  `json:"startTime"     description:"task_start_time"`                                   // task_start_time
	FinishTime     int64  `json:"finishTime"    description:"task_finish_time"`                                  // task_finish_time
	TaskCost       int    `json:"taskCost"      description:"task cost time(second)"`                            // task cost time(second)
	FailureReason  string `json:"failureReason"    description:"reason of failure"`                              // reason of failure
	TaskType       int    `json:"taskType"      description:"type，0-download，1-upload"`                          // type，0-download，1-upload
	SuccessCount   int64  `json:"successCount"  description:"success_count"`                                     // success_count
	UploadFileUrl  string `json:"uploadFileUrl" description:"the file url of upload type task"`                  // the file url of upload type task
	CreateTime     int64  `json:"createTime"     description:"create utc time"`                                  // create utc time
	LastUpdateTime int64  `json:"lastUpdateTime" description:"last update utc time"`                             // last update utc time
	Format         string `json:"format"    description:"format of file"`                                        // format of file
}

func SimplifyMerchantBatchTask(one *entity.MerchantBatchTask) *MerchantBatchTaskSimplify {
	if one == nil {
		return nil
	}
	if len(one.Format) == 0 {
		if len(one.DownloadUrl) > 0 {
			extension, _ := getFileExtensionFromURL(one.DownloadUrl)
			if strings.Contains(extension, "csv") {
				one.Format = "csv"
			} else {
				one.Format = "xlsx"
			}
		} else {
			one.Format = "xlsx"
		}
	}
	return &MerchantBatchTaskSimplify{
		Id:             one.Id,
		MerchantId:     one.MerchantId,
		MemberId:       one.MemberId,
		TaskName:       one.TaskName,
		Payload:        one.Payload,
		DownloadUrl:    link.GetTaskDownloadUrl(one),
		Status:         one.Status,
		StartTime:      one.StartTime,
		FinishTime:     one.FinishTime,
		TaskCost:       one.TaskCost,
		FailureReason:  one.FailReason,
		TaskType:       one.TaskType,
		SuccessCount:   one.SuccessCount,
		UploadFileUrl:  one.UploadFileUrl,
		CreateTime:     one.CreateTime,
		LastUpdateTime: one.LastUpdateTime,
		Format:         one.Format,
	}
}

func getFileExtensionFromURL(downloadURL string) (string, error) {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}
	filePath := parsedURL.Path
	extension := filepath.Ext(filePath)

	return extension, nil
}
