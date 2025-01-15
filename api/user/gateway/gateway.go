package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta `path:"/list" tags:"User-Gateway" method:"get" summary:"Query Gateway List"`
}
type ListRes struct {
	Gateways []*detail.Gateway `json:"gateways"`
}
