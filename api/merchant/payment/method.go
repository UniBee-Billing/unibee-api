package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
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

type MethodGetReq struct {
	g.Meta          `path:"/method_get" tags:"Payment" method:"get" summary:"Query Payment Method"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
	UserId          uint64 `json:"userId" dc:"UserId"  v:"required" `
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId"  v:"required" `
}

type MethodGetRes struct {
	Method *bean.PaymentMethod `json:"method" dc:"Method" `
}

type MethodNewReq struct {
	g.Meta         `path:"/method_new" tags:"Payment" method:"post" summary:"Create New Payment Method And Attach To User"`
	UserId         uint64      `json:"userId" dc:"UserId"   v:"required" `
	GatewayId      uint64      `json:"gatewayId" dc:"GatewayId"   v:"required" `
	Currency       string      `json:"currency" dc:""  v:"required" `
	SubscriptionId string      `json:"subscriptionId" dc:"" dc:"if provide, bind to it"`
	RedirectUrl    string      `json:"redirectUrl" dc:"Redirect Url"`
	Type           string      `json:"type"`
	Data           *gjson.Json `json:"data"`
}

type MethodNewRes struct {
	Url    string              `json:"url" dc:"Url" `
	Method *bean.PaymentMethod `json:"method" dc:"Method" `
}
