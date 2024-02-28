package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Merchant-Gateway-Controller" method:"get" summary:"Gateway List"`
}
type ListRes struct {
	Gateways []*ro.GatewaySimplify `json:"gateways"`
}

type SetupReq struct {
	g.Meta        `path:"/setup" tags:"Merchant-Gateway-Controller" method:"post" summary:"Gateway Setup"`
	GatewayName   string `p:"gatewayName"  dc:"GatewayName, stripe|paypal" v:"required"`
	GatewayKey    string `p:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `p:"gatewaySecret"  dc:"GatewaySecret" `
}
type SetupRes struct {
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Merchant-Gateway-Controller" method:"post" summary:"Gateway Webhook Edit"`
	GatewayId     uint64 `p:"gatewayId"  dc:"GatewayId" v:"required"`
	GatewayKey    string `p:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `p:"gatewaySecret"  dc:"GatewaySecret" `
}
type EditRes struct {
}

type SetupGatewayWebhookReq struct {
	g.Meta    `path:"/setup_webhook" tags:"Merchant-Gateway-Controller" method:"post" summary:"Gateway Webhook Setup"`
	GatewayId uint64 `p:"gatewayId"  dc:"GatewayId" v:"required"`
}
type SetupGatewayWebhookRes struct {
}
