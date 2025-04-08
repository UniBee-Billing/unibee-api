package gateway

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
)

type ListReq struct {
	g.Meta     `path:"/list" tags:"Checkout" method:"get" summary:"Query Gateway List"`
	MerchantId uint64 `json:"merchantId" description:"" v:"required"`
}
type ListRes struct {
	Gateways []*detail.Gateway `json:"gateways"`
}
