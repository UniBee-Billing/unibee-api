package mock

import "github.com/gogf/gf/v2/frame/g"

type CaptureReq struct {
	g.Meta     `path:"/capture" tags:"Open-Mock-Controller" method:"post" summary:"Mock Capture Payment"`
	PaymentId  string `p:"paymentId" dc:"PaymentId" v:"required"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	Currency   string `p:"currency" dc:"Currency" v:"required"`
	Amount     int64  `p:"amount" dc:" Amount, Cent" v:"required"`
}
type CaptureRes struct {
}
