package payment

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type PaymentCallbackAgainReq struct {
	g.Meta    `path:"/payment_callback_again" tags:"System-Admin" method:"post" summary:"Admin Trigger Payment Callback"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required#Require paymentId"`
}

type PaymentCallbackAgainRes struct {
}

type PaymentGatewayDetailReq struct {
	g.Meta    `path:"/payment_gateway_detail" tags:"System-Admin" method:"post" summary:"Admin Trigger Payment Callback"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required#Require paymentId"`
}

type PaymentGatewayDetailRes struct {
	PaymentDetail *gjson.Json `json:"paymentDetail"`
}
