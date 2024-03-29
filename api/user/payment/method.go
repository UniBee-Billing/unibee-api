package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type MethodListReq struct {
	g.Meta    `path:"/method_list" tags:"User-Payment-Method" method:"get" summary:"Query User Payment Method List"`
	GatewayId uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
	PaymentId string `json:"paymentId" dc:"PaymentId"  `
}

type MethodListRes struct {
	MethodList []*bean.PaymentMethod `json:"methodList" dc:"MethodList" `
}

type MethodNewReq struct {
	g.Meta    `path:"/method_new" tags:"User-Payment-Method" method:"post" summary:"User Create New Payment Method"`
	GatewayId uint64      `json:"gatewayId" dc:"GatewayId"   v:"required" `
	Type      string      `json:"type"`
	Data      *gjson.Json `json:"data"`
}

type MethodNewRes struct {
	Method *bean.PaymentMethod `json:"method" dc:"Method" `
}
