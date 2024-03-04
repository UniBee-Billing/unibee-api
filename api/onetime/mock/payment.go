package mock

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type SamplePaymentNetherlandsReq struct {
	g.Meta      `path:"/quick_pay_sample_netherlands" tags:"Open-Mock" method:"post" summary:"Mock Netherlands Create Payment (Auto Fill) https://docs.klarna.com/resources/test-environment/sample-customer-data/#netherlands"`
	Currency    string `json:"currency" dc:"Currency" v:"required"`
	Amount      int64  `json:"amount" dc:" Amount, Cent" v:"required"`
	GatewayName string `json:"gatewayName" dc:"Gateway，klarna_paynow|klarna|klarna_account|paypal" v:"required"`
	ReturnUrl   string `json:"returnUrl" dc:"Return Url" v:""`
}
type SamplePaymentNetherlandsRes struct {
	Status            string      `json:"status" dc:"Status"`
	PaymentId         string      `json:"paymentId" dc:"PaymentId"`
	MerchantPaymentId string      `json:"merchantPaymentId" dc:"MerchantPaymentId"`
	Action            *gjson.Json `json:"action" dc:"action"`
}

type DetailPayReq struct {
	g.Meta    `path:"/detail_pay" tags:"Open-Mock" method:"post" summary:"Mock Payment Detail"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
}
type DetailPayRes struct {
}
