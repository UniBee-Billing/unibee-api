package email

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TemplateListReq struct {
	g.Meta `path:"/template_list" tags:"EmailTemplate" method:"get" summary:"EmailTemplateList"`
}

type TemplateListRes struct {
	EmailTemplateList []*bean.MerchantEmailTemplateSimplify `json:"emailTemplateList" description:"Email Template Object List" `
	Total             int                                   `json:"total" dc:"Total"`
}

type TemplateUpdateReq struct {
	g.Meta          `path:"/template_update" tags:"EmailTemplate" method:"post" summary:"EmailTemplateUpdate" dc:"Update the email template"`
	TemplateName    string `json:"templateName" dc:"The name of email template"       v:"required"`
	TemplateTitle   string `json:"templateTitle" dc:"The title of email template"      v:"required"`
	TemplateContent string `json:"templateContent" dc:"The content of email template"    v:"required"`
}

type TemplateUpdateRes struct {
}

type TemplateSetDefaultReq struct {
	g.Meta       `path:"/template_set_default" tags:"EmailTemplate" method:"post" summary:"EmailTemplateSetDefault" dc:"Setup email template as default"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateSetDefaultRes struct {
}

type TemplateActivateReq struct {
	g.Meta       `path:"/template_activate" tags:"EmailTemplate" method:"post" summary:"EmailTemplateActivate" dc:"Activate the email template"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateActivateRes struct {
}

type TemplateDeactivateReq struct {
	g.Meta       `path:"/template_deactivate" tags:"EmailTemplate" method:"post" summary:"EmailTemplateDeactivate" dc:"Deactivate the email template"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateDeactivateRes struct {
}
