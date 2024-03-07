package email

import "github.com/gogf/gf/v2/frame/g"

type GatewaySetupReq struct {
	g.Meta      `path:"/gateway_setup" tags:"Email" method:"post" summary:"Merchant Email Gateway Setup"`
	GatewayName string `json:"gatewayName"  dc:"GatewayName, e.m. sendgrid" v:"required"`
	Data        string `json:"data" dc:"data" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}

type GatewaySetupRes struct {
}
