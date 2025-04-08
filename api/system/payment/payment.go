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

type PaymentGatewayCheckReq struct {
	g.Meta    `path:"/payment_gateway_checker" tags:"System-Admin" method:"post" summary:"Admin Trigger Payment Checker"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required#Require paymentId"`
}

type PaymentGatewayCheckRes struct {
}

type GetPaymentExchangeRateReq struct {
	g.Meta       `path:"/get_payment_exchange_rate" tags:"System-Admin" method:"get,post" summary:"Get Cloud Exchange Rate"`
	FromCurrency string `json:"fromCurrency" dc:"From Currency"`
	ToCurrency   string `json:"toCurrency" dc:"To Currency"`
}

type GetPaymentExchangeRateRes struct {
	ExchangeRate float64 `json:"exchangeRate"`
}
