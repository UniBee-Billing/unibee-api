// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantEmailTemplate is the golang structure for table merchant_email_template.
type MerchantEmailTemplate struct {
	Id                 int64       `json:"id"                 description:""`                      //
	MerchantId         uint64      `json:"merchantId"         description:""`                      //
	TemplateName       string      `json:"templateName"       description:""`                      //
	TemplateTitle      string      `json:"templateTitle"      description:""`                      //
	TemplateContent    string      `json:"templateContent"    description:""`                      //
	TemplateAttachName string      `json:"templateAttachName" description:""`                      //
	GmtCreate          *gtime.Time `json:"gmtCreate"          description:"create time"`           // create time
	GmtModify          *gtime.Time `json:"gmtModify"          description:"update time"`           // update time
	IsDeleted          int         `json:"isDeleted"          description:"0-UnDeleted，1-Deleted"` // 0-UnDeleted，1-Deleted
	CreateTime         int64       `json:"createTime"         description:"create utc time"`       // create utc time
	Status             int64       `json:"status"             description:"0-Active,1-InActive"`   // 0-Active,1-InActive
	GatewayTemplateId  string      `json:"gatewayTemplateId"  description:""`                      //
	LanguageData       string      `json:"languageData"       description:""`                      //
}
