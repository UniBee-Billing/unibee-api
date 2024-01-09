// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileUpload is the golang structure for table file_upload.
type FileUpload struct {
	Id        int64       `json:"id"        ` //
	UserId    string      `json:"userId"    ` //
	Url       string      `json:"url"       ` //
	FileName  string      `json:"fileName"  ` //
	Tag       string      `json:"tag"       ` //
	GmtCreate *gtime.Time `json:"gmtCreate" ` //
	GmtModify *gtime.Time `json:"gmtModify" ` //
	IsDeleted int         `json:"isDeleted" ` // 是否删除，0-未删除，1-删除
}
