package email

import "github.com/gogf/gf/v2/frame/g"

type GatewaySetupReq struct {
	g.Meta      `path:"/gateway_setup" tags:"Email" method:"post" summary:"EmailGatewaySetup"`
	GatewayName string `json:"gatewayName"  dc:"The name of email gateway, 'sendgrid' or other for future updates" v:"required"`
	Data        string `json:"data" dc:"The setup data of email gateway" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"Whether setup the gateway as default or not, default is true" `
}

type GatewaySetupRes struct {
}

type SendTemplateEmailToUserReq struct {
	g.Meta       `path:"/send_template_email_to_user" tags:"Email" method:"post" summary:"SendTemplateEmailToUser"`
	TemplateName string                 `json:"templateName" dc:"The name of email template"       v:"required"`
	UserId       int64                  `json:"userId" dc:"UserId" v:"required" `
	Variables    map[string]interface{} `json:"variables" dc:"Variables，Map"`
}

type SendTemplateEmailToUserRes struct {
}
