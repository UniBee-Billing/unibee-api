// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmailTemplate is the golang structure for table email_template.
type EmailTemplate struct {
	Id                 int64       `json:"id"                 description:""`                      //
	MerchantId         int64       `json:"merchantId"         description:""`                      //
	TemplateName       string      `json:"templateName"       description:""`                      //
	TemplateTitle      string      `json:"templateTitle"      description:""`                      //
	TemplateContent    string      `json:"templateContent"    description:""`                      //
	TemplateAttachName string      `json:"templateAttachName" description:""`                      //
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`           // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`           // update time
	IsDeleted          int         `json:"isDeleted"          description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
}
