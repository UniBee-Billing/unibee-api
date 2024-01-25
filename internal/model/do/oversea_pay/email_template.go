// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmailTemplate is the golang structure of table email_template for DAO operations like Where/Data.
type EmailTemplate struct {
	g.Meta             `orm:"table:email_template, do:true"`
	Id                 interface{} //
	MerchantId         interface{} //
	TemplateName       interface{} //
	TemplateTitle      interface{} //
	TemplateContent    interface{} //
	TemplateAttachName interface{} //
	GmtCreate          *gtime.Time //
	GmtModify          *gtime.Time //
	IsDeleted          interface{} // 是否删除，0-未删除，1-删除
}
