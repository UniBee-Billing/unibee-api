// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FileUpload is the golang structure of table file_upload for DAO operations like Where/Data.
type FileUpload struct {
	g.Meta    `orm:"table:file_upload, do:true"`
	Id        interface{} //
	UserId    interface{} //
	Url       interface{} //
	FileName  interface{} //
	Tag       interface{} //
	GmtCreate *gtime.Time // create time
	GmtModify *gtime.Time //
	IsDeleted interface{} // 0-UnDeleted，1-Deleted
}
