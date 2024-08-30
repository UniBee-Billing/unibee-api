// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileUpload is the golang structure for table file_upload.
type FileUpload struct {
	Id         int64       `json:"id"         description:""`                      //
	UserId     string      `json:"userId"     description:""`                      //
	Url        string      `json:"url"        description:""`                      //
	FileName   string      `json:"fileName"   description:""`                      //
	Tag        string      `json:"tag"        description:""`                      //
	GmtCreate  *gtime.Time `json:"gmtCreate"  description:"create time"`           // create time
	GmtModify  *gtime.Time `json:"gmtModify"  description:""`                      //
	IsDeleted  int         `json:"isDeleted"  description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime int64       `json:"createTime" description:"create utc time"`       // create utc time
	Data       []byte      `json:"data"       description:""`                      //
}
