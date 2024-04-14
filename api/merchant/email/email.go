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
