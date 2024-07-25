// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MerchantEmailTemplate is the golang structure of table merchant_email_template for DAO operations like Where/Data.
type MerchantEmailTemplate struct {
	g.Meta             `orm:"table:merchant_email_template, do:true"`
	Id                 interface{} //
	MerchantId         interface{} //
	TemplateName       interface{} //
	TemplateTitle      interface{} //
	TemplateContent    interface{} //
	TemplateAttachName interface{} //
	GmtCreate          *gtime.Time // create time
	GmtModify          *gtime.Time // update time
	IsDeleted          interface{} // 0-UnDeleted，1-Deleted
	CreateTime         interface{} // create utc time
	Status             interface{} // 0-Active,1-InActive
}
