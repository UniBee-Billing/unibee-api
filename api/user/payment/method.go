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

type MethodGetReq struct {
	g.Meta          `path:"/method_get" tags:"User-Payment-Method" method:"get" summary:"Query Payment Method"`
	GatewayId       uint64 `json:"gatewayId" dc:"GatewayId"   v:"required" `
	PaymentMethodId string `json:"paymentMethodId" dc:"PaymentMethodId"  v:"required" `
}

type MethodGetRes struct {
	Method *bean.PaymentMethod `json:"method" dc:"Method" `
}

type MethodNewReq struct {
	g.Meta         `path:"/method_new" tags:"User-Payment-Method" method:"post" summary:"User Create New Payment Method"`
	GatewayId      uint64      `json:"gatewayId" dc:"GatewayId"   v:"required" `
	Currency       string      `json:"currency" dc:""  v:"required" `
	SubscriptionId string      `json:"subscriptionId" dc:"if provide, bind to it"`
	RedirectUrl    string      `json:"redirectUrl" dc:"Redirect Url"`
	Type           string      `json:"type" dc:""`
	Data           *gjson.Json `json:"data" dc:""`
}

type MethodNewRes struct {
	Url    string              `json:"url" dc:"Url" `
	Method *bean.PaymentMethod `json:"method" dc:"Method" `
}
