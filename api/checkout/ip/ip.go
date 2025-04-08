package ip

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type ResolveReq struct {
	g.Meta `path:"/resolve" tags:"Checkout" method:"get" summary:"Checkout IP Resolve"`
	IP     string `json:"ip" dc:"ip" v:"required"`
}
type ResolveRes struct {
	Location *gjson.Json `json:"location" dc:"Ip Location Data"`
}
