package email

import "github.com/gogf/gf/v2/frame/g"

type GatewaySetupReq struct {
	g.Meta      `path:"/gateway_setup" tags:"Email" method:"post" summary:"Email Gateway Setup"`
	GatewayName string `json:"gatewayName"  dc:"The name of email gateway, 'sendgrid' or other for future updates" v:"required"`
	Data        string `json:"data" dc:"The setup data of email gateway" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"Whether setup the gateway as default or not, default is true" `
}

type GatewaySetupRes struct {
}

type SendTemplateEmailToUserReq struct {
	g.Meta       `path:"/send_template_email_to_user" tags:"Email" method:"post" summary:"Send Template Email To User"`
	TemplateName string                 `json:"templateName" dc:"The name of email template"       v:"required"`
	UserId       int64                  `json:"userId" dc:"UserId" v:"required" `
	Variables    map[string]interface{} `json:"variables" dc:"Variables，Map"`
}

type SendTemplateEmailToUserRes struct {
}

type SenderSetupReq struct {
	g.Meta  `path:"/email_sender_setup" tags:"Email" method:"post" summary:"Email Sender Setup"`
	Name    string `json:"name"  dc:"The name of email sender, like 'no-reply'" v:"required"`
	Address string `json:"address" dc:"The address of email sender, like 'no-reply@unibee.dev'" v:"required"`
}

type SenderSetupRes struct {
}
