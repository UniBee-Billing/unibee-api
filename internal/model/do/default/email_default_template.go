// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmailDefaultTemplate is the golang structure of table email_default_template for DAO operations like Where/Data.
type EmailDefaultTemplate struct {
	g.Meta              `orm:"table:email_default_template, do:true"`
	Id                  interface{} //
	TemplateName        interface{} //
	TemplateDescription interface{} //
	TemplateTitle       interface{} //
	TemplateContent     interface{} //
	TemplateAttachName  interface{} //
	GmtCreate           *gtime.Time // create time
	GmtModify           *gtime.Time // update time
	IsDeleted           interface{} // 0-UnDeleted，1-Deleted
	CreateTime          interface{} // create utc time
}
