package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type MethodListReq struct {
	g.Meta    `path:"/method_list" tags:"User-Payment-Method" method:"get" summary:"Query User Payment Method List"`
	GatewayId uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
	UserId    uint64 `json:"userId" dc:"UserId"   v:"required" `
	PaymentId string `json:"paymentId" dc:"PaymentId"  `
}

type MethodListRes struct {
	MethodList []*ro.PaymentMethod `json:"methodList" dc:"MethodList" `
}
