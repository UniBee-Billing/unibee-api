package mock

import "github.com/gogf/gf/v2/frame/g"

type CaptureReq struct {
	g.Meta    `path:"/capture" tags:"Open-Mock" method:"post" summary:"Mock Capture Payment"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
	Currency  string `json:"currency" dc:"Currency" v:"required"`
	Amount    int64  `json:"amount" dc:" Amount, Cent" v:"required"`
}
type CaptureRes struct {
}
