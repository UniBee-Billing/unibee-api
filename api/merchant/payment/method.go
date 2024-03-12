package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type MethodListReq struct {
	g.Meta    `path:"/method_list" tags:"Payment" method:"get" summary:"Query Payment Method List"`
	GatewayId uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
	UserId    uint64 `json:"userId" dc:"UserId" `
	PaymentId string `json:"paymentId" dc:"PaymentId"  `
}

type MethodListRes struct {
	MethodList []*bean.PaymentMethod `json:"methodList" dc:"MethodList" `
}
