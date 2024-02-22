package mock

import "github.com/gogf/gf/v2/frame/g"

type RefundReq struct {
	g.Meta     `path:"/refund" tags:"Open-Mock-Controller" method:"post" summary:"Mock Refund"`
	PaymentId  string `p:"paymentId" dc:"PaymentId" v:"required"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
	Currency   string `p:"currency" dc:"Currency" v:"required"`
	Amount     int64  `p:"amount" dc:" Amount, Cent" v:"required"`
}
type RefundRes struct {
}
