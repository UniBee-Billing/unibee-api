package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type MethodListReq struct {
	g.Meta    `path:"/method_list" tags:"User-Payment-Method" method:"get" summary:"Query User Payment Method List"`
	GatewayId uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
}

type MethodListRes struct {
	MethodList []*ro.PaymentMethod `json:"methodList" dc:"MethodList" `
}

type NewReq struct {
	g.Meta    `path:"/new" tags:"User-Payment-Method" method:"post" summary:"User Create New Payment Method"`
	GatewayId uint64      `json:"gatewayId" dc:"GatewayId"   v:"required" `
	Type      string      `json:"type"`
	Data      *gjson.Json `json:"data"`
}

type NewRes struct {
	Method *ro.PaymentMethod `json:"method" dc:"Method" `
}
