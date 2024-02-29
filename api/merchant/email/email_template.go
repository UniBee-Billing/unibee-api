package email

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/query"
)

type MerchantEmailTemplateListReq struct {
	g.Meta `path:"/merchant_email_template_list" tags:"Merchant-Email-Template-Controller" method:"get" summary:"Merchant Email Template List"`
}

type MerchantEmailTemplateListRes struct {
	EmailTemplateList []*query.EmailTemplateVo `json:"emailTemplateList" description:"EmailTemplateList" `
}

type MerchantEmailTemplateUpdateReq struct {
	g.Meta          `path:"/merchant_email_template_update" tags:"Merchant-Email-Template-Controller" method:"post" summary:"Merchant Email Template Update"`
	TemplateName    string `p:"templateName" dc:"templateName"       v:"required"`
	TemplateTitle   string `p:"templateTitle" dc:"templateTitle"      v:"required"`
	TemplateContent string `p:"templateContent" dc:"templateContent"    v:"required"`
}

type MerchantEmailTemplateUpdateRes struct {
}

type MerchantEmailTemplateSetDefaultReq struct {
	g.Meta       `path:"/merchant_email_template_set_default" tags:"Merchant-Email-Template-Controller" method:"post" summary:"Merchant Email Template Set Default"`
	TemplateName string `p:"templateName" dc:"templateName" v:"required"`
}

type MerchantEmailTemplateSetDefaultRes struct {
}

type MerchantEmailTemplateActivateReq struct {
	g.Meta       `path:"/merchant_email_template_activate" tags:"Merchant-Email-Template-Controller" method:"post" summary:"Merchant Email Template Activate"`
	TemplateName string `p:"templateName" dc:"templateName" v:"required"`
}

type MerchantEmailTemplateActivateRes struct {
}

type MerchantEmailTemplateDeactivateReq struct {
	g.Meta       `path:"/merchant_email_template_deactivate" tags:"Merchant-Email-Template-Controller" method:"post" summary:"Merchant Email Template Deactivate"`
	TemplateName string `p:"templateName" dc:"templateName" v:"required"`
}

type MerchantEmailTemplateDeactivateRes struct {
}

type MerchantEmailGatewaySetupReq struct {
	g.Meta      `path:"/email_gateway_setup" tags:"Merchant-Email-Template-Controller" method:"post" summary:"Merchant Email Gateway Setup"`
	GatewayName string `p:"gatewayName"  dc:"GatewayName, e.m. sendgrid" v:"required"`
	Data        string `p:"data" dc:"data" v:"required"`
	IsDefault   bool   `p:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}

type MerchantEmailGatewaySetupRes struct {
}
