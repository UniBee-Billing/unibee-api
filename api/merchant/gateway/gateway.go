package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Gateway" method:"get" summary:"PaymentGatewayList"`
}
type ListRes struct {
	Gateways []*bean.GatewaySimplify `json:"gateways"`
}

type SetupReq struct {
	g.Meta        `path:"/setup" tags:"Gateway" method:"post" summary:"PaymentGatewaySetup"`
	GatewayName   string `json:"gatewayName"  dc:"GatewayName, stripe|paypal|changelly" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `json:"gatewaySecret"  dc:"GatewaySecret" `
}
type SetupRes struct {
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Gateway" method:"post" summary:"PaymentGatewayEdit"`
	GatewayId     uint64 `json:"gatewayId"  dc:"GatewayId" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `json:"gatewaySecret"  dc:"GatewaySecret" `
}
type EditRes struct {
}

type EditCountryConfigReq struct {
	g.Meta        `path:"/edit_country_config" tags:"Gateway" method:"post" summary:"PaymentGatewayCountryConfigEdit"`
	GatewayId     uint64          `json:"gatewayId"  dc:"GatewayId" v:"required"`
	CountryConfig map[string]bool `json:"countryConfig"  dc:"CountryConfig" `
}
type EditCountryConfigRes struct {
}

type SetupWebhookReq struct {
	g.Meta    `path:"/setup_webhook" tags:"Gateway" method:"post" summary:"PaymentGatewayWebhookSetup"`
	GatewayId uint64 `json:"gatewayId"  dc:"GatewayId" v:"required"`
}
type SetupWebhookRes struct {
	GatewayWebhookUrl string `json:"gatewayWebhookUrl"  dc:"GatewayWebhookUrl"`
}
