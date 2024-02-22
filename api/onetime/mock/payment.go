package mock

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type SamplePaymentNetherlandsReq struct {
	g.Meta      `path:"/quick_pay_sample_netherlands" tags:"Open-Mock-Controller" method:"post" summary:"Mock Netherlands Create Payment (Auto Fill) https://docs.klarna.com/resources/test-environment/sample-customer-data/#netherlands"`
	Currency    string `p:"currency" dc:"Currency" v:"required"`
	Amount      int64  `p:"amount" dc:" Amount, Cent" v:"required"`
	GatewayName string `p:"gatewayName" dc:"Gateway，klarna_paynow|klarna|klarna_account|paypal" v:"required"`
	ReturnUrl   string `p:"returnUrl" dc:"Return Url" v:""`
}
type SamplePaymentNetherlandsRes struct {
	Status            string      `p:"status" dc:"Status"`
	PaymentId         string      `p:"paymentId" dc:"PaymentId"`
	MerchantPaymentId string      `p:"merchantPaymentId" dc:"MerchantPaymentId"`
	Action            *gjson.Json `p:"action" dc:"action"`
}

type DetailPayReq struct {
	g.Meta    `path:"/detail_pay" tags:"Open-Mock-Controller" method:"post" summary:"Mock Payment Detail"`
	PaymentId string `p:"paymentId" dc:"PaymentId" v:"required"`
}
type DetailPayRes struct {
}
