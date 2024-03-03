package mock

import "github.com/gogf/gf/v2/frame/g"

type RefundReq struct {
	g.Meta    `path:"/refund" tags:"Open-Mock-Controller" method:"post" summary:"Mock Refund"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
	Currency  string `json:"currency" dc:"Currency" v:"required"`
	Amount    int64  `json:"amount" dc:" Amount, Cent" v:"required"`
}
type RefundRes struct {
}
