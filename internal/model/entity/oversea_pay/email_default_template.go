// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmailDefaultTemplate is the golang structure for table email_default_template.
type EmailDefaultTemplate struct {
	Id                  int64       `json:"id"                  description:""`                      //
	TemplateName        string      `json:"templateName"        description:""`                      //
	TemplateDescription string      `json:"templateDescription" description:""`                      //
	TemplateTitle       string      `json:"templateTitle"       description:""`                      //
	TemplateContent     string      `json:"templateContent"     description:""`                      //
	TemplateAttachName  string      `json:"templateAttachName"  description:""`                      //
	GmtCreate           *gtime.Time `json:"gmtCreate"           description:"create time"`           // create time
	GmtModify           *gtime.Time `json:"gmtModify"           description:"update time"`           // update time
	IsDeleted           int         `json:"isDeleted"           description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime          int64       `json:"createTime"          description:"create utc time"`       // create utc time
}
