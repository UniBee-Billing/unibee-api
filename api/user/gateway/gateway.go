package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"User-Gateway" method:"get" summary:"Query Gateway List"`
}
type ListRes struct {
	Gateways []*ro.GatewaySimplify `json:"gateways"`
}
