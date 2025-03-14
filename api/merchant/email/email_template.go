package email

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type TemplateListReq struct {
	g.Meta `path:"/template_list" tags:"Email Template" method:"get" summary:"Get Email Template List"`
}

type TemplateListRes struct {
	EmailTemplateList []*bean.MerchantEmailTemplate `json:"emailTemplateList" description:"Email Template Object List" `
	Total             int                           `json:"total" dc:"Total"`
}

type TemplateUpdateReq struct {
	g.Meta          `path:"/template_update" tags:"Email Template" method:"post" summary:"Email Template Update" dc:"Update the email template"`
	TemplateName    string `json:"templateName" dc:"The name of email template"       v:"required"`
	TemplateTitle   string `json:"templateTitle" dc:"The title of email template"      v:"required"`
	TemplateContent string `json:"templateContent" dc:"The content of email template"    v:"required"`
}

type TemplateUpdateRes struct {
}

type TemplateSetDefaultReq struct {
	g.Meta       `path:"/template_set_default" tags:"Email Template" method:"post" summary:"Setup Email Template Default" dc:"Setup email template as default"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateSetDefaultRes struct {
}

type TemplateActivateReq struct {
	g.Meta       `path:"/template_activate" tags:"Email Template" method:"post" summary:"Email Template Activate" dc:"Activate the email template"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateActivateRes struct {
}

type TemplateDeactivateReq struct {
	g.Meta       `path:"/template_deactivate" tags:"Email Template" method:"post" summary:"Email Template Deactivate" dc:"Deactivate the email template"`
	TemplateName string `json:"templateName" dc:"The name of email template" v:"required"`
}

type TemplateDeactivateRes struct {
}

type CustomizeLocalizationTemplateSyncReq struct {
	g.Meta        `path:"/custom_localization_template_sync" tags:"Email Template" method:"post" summary:"Customize Localization Template Sync" dc:"Sync the custom localization email template to gateway (sendgrid)"`
	TemplateData  []bean.MerchantLocalizationEmailTemplate `json:"templateData" dc:"TemplateData" v:"required"`
	VersionEnable bool                                     `json:"versionEnable" dc:"VersionEnable"`
}

type CustomizeLocalizationTemplateSyncRes struct {
}
