package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type EditSortReq struct {
	g.Meta       `path:"/edit_sort" tags:"Gateway" method:"post" summary:"Edit Payment Gateway Sort"`
	GatewaySorts []*detail.GatewaySort `json:"gatewaySorts" dc:"structure gateway sort"`
}

type EditSortRes struct {
	Gateways []*detail.Gateway `json:"gateways" dc:"Payment Gateway Setup Object List"`
}

type SetupListReq struct {
	g.Meta `path:"/setup_list" tags:"Gateway" method:"get" summary:"Get Payment Gateway Setup List"`
}
type SetupListRes struct {
	Gateways []*detail.Gateway `json:"gateways" dc:"Payment Gateway Setup Object List"`
}

type DetailReq struct {
	g.Meta      `path:"/detail" tags:"Gateway" method:"get,post" summary:"Payment Gateway" dc:"Get Payment Gateway Detail"`
	GatewayId   *uint64 `json:"gatewayId"  dc:"The id of payment gateway, either gatewayId or gatewayName"`
	GatewayName *string `json:"gatewayName"  dc:"The name of payment gateway, , either gatewayId or gatewayName, stripe|paypal|changelly|unitpay|payssion|cryptadium"`
}
type DetailRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type ListReq struct {
	g.Meta  `path:"/list" tags:"Gateway" method:"get" summary:"Get Payment Gateway List"`
	Archive *bool `json:"archive" dc:"Filter archive gateway or not, default all"`
}
type ListRes struct {
	Gateways []*detail.Gateway `json:"gateways" dc:"Payment Gateway Object List"`
}

type SetupReq struct {
	g.Meta              `path:"/setup" tags:"Gateway" method:"post" summary:"Payment Gateway Setup" dc:"Setup Payment gateway"`
	GatewayName         string                            `json:"gatewayName"  dc:"The name of payment gateway, stripe|paypal|changelly|unitpay|payssion|cryptadium" v:"required"`
	DisplayName         *string                           `json:"displayName"  dc:"The displayName of payment gateway"`
	GatewayIcons        *[]string                         `json:"gatewayIcons"  dc:"The icons of payment gateway"`
	Sort                *int64                            `json:"sort"  dc:"The sort value of payment gateway, The bigger, the closer to the front"`
	GatewayKey          string                            `json:"gatewayKey"  dc:"The key of payment gateway" `
	GatewaySecret       string                            `json:"gatewaySecret"  dc:"The secret of payment gateway" `
	SubGateway          string                            `json:"subGateway"  dc:"The sub gateway of payment gateway" `
	CurrencyExchange    []*detail.GatewayCurrencyExchange `json:"currencyExchange" dc:"The currency exchange for gateway payment, effect at start of payment creation when currency matched"`
	GatewayPaymentTypes []string                          `json:"gatewayPaymentTypes"  dc:"Selected gateway payment types"`
}
type SetupRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type EditReq struct {
	g.Meta              `path:"/edit" tags:"Gateway" method:"post" summary:"Payment Gateway Edit" dc:"Edit the exist payment gateway"`
	GatewayId           uint64                            `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	DisplayName         *string                           `json:"displayName"  dc:"The displayName of payment gateway"`
	GatewayLogo         *[]string                         `json:"gatewayLogo"  dc:"The logo of payment gateway"`
	Sort                *int64                            `json:"sort"  dc:"The sort value of payment gateway, The bigger, the closer to the front"`
	GatewayKey          *string                           `json:"gatewayKey"  dc:"The key of payment gateway" `
	GatewaySecret       *string                           `json:"gatewaySecret"  dc:"The secret of payment gateway" `
	SubGateway          *string                           `json:"subGateway"  dc:"The sub gateway of payment gateway" `
	CurrencyExchange    []*detail.GatewayCurrencyExchange `json:"currencyExchange" dc:"The currency exchange for gateway payment, effect at start of payment creation when currency matched"`
	GatewayPaymentTypes []string                          `json:"gatewayPaymentTypes"  dc:"Selected gateway payment types"`
}
type EditRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type ArchiveReq struct {
	g.Meta    `path:"/archive" tags:"Gateway" method:"post" summary:"Payment Gateway Archive" dc:"Archive the exist payment gateway，available for setup new one, the old user continue using the old gateway"`
	GatewayId uint64 `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
}
type ArchiveRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type EditCountryConfigReq struct {
	g.Meta        `path:"/edit_country_config" tags:"Gateway" method:"post" summary:"Payment Gateway Country Config Edit" dc:"Edit country config for payment gateway, to enable or disable the payment for countryCode, default is enable"`
	GatewayId     uint64          `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	CountryConfig map[string]bool `json:"countryConfig"  dc:"The country config of payment gateway, a map with countryCode as key, and value for enable or disable" `
}
type EditCountryConfigRes struct {
}

type SetupWebhookReq struct {
	g.Meta        `path:"/setup_webhook" tags:"Gateway" method:"post" summary:"Payment Gateway Webhook Setup"`
	GatewayId     uint64 `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	WebhookSecret string `json:"webhookSecret"  dc:"The secret of gateway webhook"`
}
type SetupWebhookRes struct {
	GatewayWebhookUrl string `json:"gatewayWebhookUrl"  dc:"The webhook endpoint url of payment gateway, if gateway is stripe, the url will setting to stripe by api automatic"`
}

type WireTransferSetupReq struct {
	g.Meta        `path:"/wire_transfer_setup" tags:"Gateway" method:"post" summary:"Wire Transfer Setup" dc:"Setup the wire transfer"`
	DisplayName   *string             `json:"displayName"  dc:"The displayName of payment gateway"`
	GatewayIcons  *[]string           `json:"gatewayIcons"  dc:"The icons of payment gateway"`
	Sort          *int64              `json:"sort"  dc:"The sort value of payment gateway, The bigger, the closer to the front"`
	Currency      string              `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64               `json:"minimumAmount"   dc:"The minimum amount of wire transfer, cents" v:"required" `
	Bank          *detail.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer" v:"required"`
}
type WireTransferSetupRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type WireTransferEditReq struct {
	g.Meta        `path:"/wire_transfer_edit" tags:"Gateway" method:"post" summary:"Wire Transfer Edit" dc:"Edit the wire transfer"`
	GatewayId     uint64              `json:"gatewayId"  dc:"The id of payment gateway" v:"required"`
	DisplayName   *string             `json:"displayName"  dc:"The displayName of payment gateway"`
	GatewayIcons  *[]string           `json:"gatewayIcons"  dc:"The icons of payment gateway"`
	Sort          *int64              `json:"sort"  dc:"The sort value of payment gateway, The bigger, the closer to the front"`
	Currency      string              `json:"currency"   dc:"The currency of wire transfer " v:"required" `
	MinimumAmount int64               `json:"minimumAmount"   dc:"The minimum amount of wire transfer, cents" v:"required" `
	Bank          *detail.GatewayBank `json:"bank"   dc:"The receiving bank of wire transfer" v:"required"`
}
type WireTransferEditRes struct {
	Gateway *detail.Gateway `json:"gateway" dc:"Payment Gateway Object"`
}

type SetupExchangeApiReq struct {
	g.Meta             `path:"/setup_exchange_rate_api" tags:"Gateway" method:"post" summary:"Exchange Rate Api Setup"`
	ExchangeRateApiKey string `json:"exchangeRateApiKey"  dc:"The key of exchange rate api"`
}
type SetupExchangeApiRes struct {
	Data string `json:"data"  dc:"The hide star key of exchange rate api"`
}
