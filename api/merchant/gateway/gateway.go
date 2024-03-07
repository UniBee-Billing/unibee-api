package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"Gateway" method:"get" summary:"Gateway List"`
}
type ListRes struct {
	Gateways []*ro.GatewaySimplify `json:"gateways"`
}

type SetupReq struct {
	g.Meta        `path:"/setup" tags:"Gateway" method:"post" summary:"Gateway Setup"`
	GatewayName   string `json:"gatewayName"  dc:"GatewayName, stripe|paypal" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `json:"gatewaySecret"  dc:"GatewaySecret" `
}
type SetupRes struct {
}

type EditReq struct {
	g.Meta        `path:"/edit" tags:"Gateway" method:"post" summary:"Gateway Webhook Edit"`
	GatewayId     uint64 `json:"gatewayId"  dc:"GatewayId" v:"required"`
	GatewayKey    string `json:"gatewayKey"  dc:"GatewayKey" `
	GatewaySecret string `json:"gatewaySecret"  dc:"GatewaySecret" `
}
type EditRes struct {
}

type SetupWebhookReq struct {
	g.Meta    `path:"/setup_webhook" tags:"Gateway" method:"post" summary:"Gateway Webhook Setup"`
	GatewayId uint64 `json:"gatewayId"  dc:"GatewayId" v:"required"`
}
type SetupWebhookRes struct {
}
