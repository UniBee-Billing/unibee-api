package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Gateway" method:"get" summary:"PaymentGatewayList"`
}
type ListRes struct {
	Gateways []*bean.Gateway `json:"gateways" dc:"Payment Gateway Object List"`
}

type SetupReq struct {
	g.Meta        `path:"/setup" tags:"Gateway" method:"post" summary:"PaymentGatewaySetup" dc:"Setup the payment gateway"`
	GatewayName   string `json:"gatewayName"  dc:"The name of payment gateway, stripe|paypal|changelly" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"The key of payment gateway" `
	GatewaySecret string `json:"gatewaySecret"  dc:"The secret of payment gateway" `
}
type SetupRes struct {
	Gateway *bean.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Gateway" method:"post" summary:"PaymentGatewayEdit" dc:"edit the exist payment gateway"`
	GatewayId     uint64 `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"The key of payment gateway" `
	GatewaySecret string `json:"gatewaySecret"  dc:"The secret of payment gateway" `
}
type EditRes struct {
	Gateway *bean.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type EditCountryConfigReq struct {
	g.Meta        `path:"/edit_country_config" tags:"Gateway" method:"post" summary:"PaymentGatewayCountryConfigEdit" dc:"Edit country config for payment gateway, to enable or disable the payment for countryCode, default is enable"`
	GatewayId     uint64          `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	CountryConfig map[string]bool `json:"countryConfig"  dc:"The country config of payment gateway, a map with countryCode as key, and value for enable or disable" `
}
type EditCountryConfigRes struct {
}

type SetupWebhookReq struct {
	g.Meta        `path:"/setup_webhook" tags:"Gateway" method:"post" summary:"PaymentGatewayWebhookSetup"`
	GatewayId     uint64 `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	WebhookSecret string `json:"webhookSecret"  dc:"The secret of gateway webhook"`
}
type SetupWebhookRes struct {
	GatewayWebhookUrl string `json:"gatewayWebhookUrl"  dc:"The webhook endpoint url of payment gateway, if gateway is stripe, the url will setting to stripe by api automatic"`
}

type WireTransferSetupReq struct {
	g.Meta        `path:"/wire_transfer_setup" tags:"Gateway" method:"post" summary:"WireTransferSetup" dc:"Setup the wire transfer"`
	Currency      string            `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64             `json:"minimumAmount"   dc:"The minimum amount of wire transfer, cents" v:"required" `
	Bank          *bean.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer" v:"required"`
}
type WireTransferSetupRes struct {
}

type WireTransferEditReq struct {
	g.Meta        `path:"/wire_transfer_edit" tags:"Gateway" method:"post" summary:"WireTransferEdit" dc:"Edit the wire transfer"`
	GatewayId     uint64            `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	Currency      string            `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64             `json:"minimumAmount"   dc:"The minimum amount of wire transfer, cents" v:"required" `
	Bank          *bean.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer" v:"required"`
}
type WireTransferEditRes struct {
}
