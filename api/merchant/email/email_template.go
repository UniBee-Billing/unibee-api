package email

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TemplateListReq struct {
	g.Meta `path:"/template_list" tags:"Email-Template" method:"get" summary:"Merchant Email Template List"`
}

type TemplateListRes struct {
	EmailTemplateList []*bean.EmailTemplateVo `json:"emailTemplateList" description:"EmailTemplateList" `
}

type TemplateUpdateReq struct {
	g.Meta          `path:"/template_update" tags:"Email-Template" method:"post" summary:"Merchant Email Template Update"`
	TemplateName    string `json:"templateName" dc:"templateName"       v:"required"`
	TemplateTitle   string `json:"templateTitle" dc:"templateTitle"      v:"required"`
	TemplateContent string `json:"templateContent" dc:"templateContent"    v:"required"`
}

type TemplateUpdateRes struct {
}

type TemplateSetDefaultReq struct {
	g.Meta       `path:"/template_set_default" tags:"Email-Template" method:"post" summary:"Merchant Email Template Set Default"`
	TemplateName string `json:"templateName" dc:"templateName" v:"required"`
}

type TemplateSetDefaultRes struct {
}

type TemplateActivateReq struct {
	g.Meta       `path:"/template_activate" tags:"Email-Template" method:"post" summary:"Merchant Email Template Activate"`
	TemplateName string `json:"templateName" dc:"templateName" v:"required"`
}

type TemplateActivateRes struct {
}

type TemplateDeactivateReq struct {
	g.Meta       `path:"/template_deactivate" tags:"Email-Template" method:"post" summary:"Merchant Email Template Deactivate"`
	TemplateName string `json:"templateName" dc:"templateName" v:"required"`
}

type TemplateDeactivateRes struct {
}
