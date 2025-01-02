package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type MethodListReq struct {
	g.Meta    `path:"/method_list" tags:"Payment" method:"get" summary:"Payment Method List" dc:"The method list of payment gateway"`
	GatewayId uint64 `json:"gatewayId" dc:"The unique id of gateway"   v:"required" `
	UserId    uint64 `json:"userId" dc:"The id of user" `
	PaymentId string `json:"paymentId" dc:"The unique id of payment"  `
}

type MethodListRes struct {
	MethodList []*bean.PaymentMethod `json:"methodList" dc:"Payment Method Object List" `
}

type MethodGetReq struct {
	g.Meta          `path:"/method_get" tags:"Payment" method:"get" summary:"Payment Method" dc:"The method of payment gateway"`
	GatewayId       uint64 `json:"gatewayId" dc:"The unique id of gateway"   v:"required" `
	UserId          uint64 `json:"userId" dc:"The customer's unique id"  v:"required" `
	PaymentMethodId string `json:"paymentMethodId" dc:"The unique id of payment method"  v:"required" `
}

type MethodGetRes struct {
	Method *bean.PaymentMethod `json:"method" dc:"Method Object" `
}

type MethodNewReq struct {
	g.Meta         `path:"/method_new" tags:"Payment" method:"post" summary:"Create New Payment Method"`
	UserId         uint64                 `json:"userId" dc:"The customer's unique id"   v:"required" `
	GatewayId      uint64                 `json:"gatewayId" dc:"The unique id of gateway"   v:"required" `
	Currency       string                 `json:"currency" dc:"The currency of payment method" `
	SubscriptionId string                 `json:"subscriptionId" dc:"The id of subscription that want to attach"`
	RedirectUrl    string                 `json:"redirectUrl" dc:"The redirect url when method created return back"`
	Type           string                 `json:"type"`
	Metadata       map[string]interface{} `json:"metadata" dc:"Metadata，Map"`
}

type MethodNewRes struct {
	Url    string              `json:"url" dc:"The gateway url to create a new method" `
	Method *bean.PaymentMethod `json:"method" dc:"Method Object" `
}

type MethodDeleteReq struct {
	g.Meta          `path:"/method_delete" tags:"Payment" method:"post" summary:"Delete Payment Method"`
	UserId          uint64 `json:"userId" dc:"The customer's unique id"   v:"required" `
	GatewayId       uint64 `json:"gatewayId" dc:"The unique id of gateway"   v:"required" `
	PaymentMethodId string `json:"paymentMethodId" dc:"The unique id of payment method"  v:"required" `
}

type MethodDeleteRes struct {
}
