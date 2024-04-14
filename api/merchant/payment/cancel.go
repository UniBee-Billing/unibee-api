package payment

import "github.com/gogf/gf/v2/frame/g"

type CancelReq struct {
	g.Meta    `path:"/cancel" tags:"Payment" method:"post" summary:"CancelPayment"`
	PaymentId string `json:"paymentId" dc:"The unique id of payment" v:"required"`
}
type CancelRes struct {
}

type RefundCancelReq struct {
	g.Meta   `path:"/refund/cancel" tags:"Payment" method:"post" summary:"CancelPaymentRefund"`
	RefundId string `json:"refundId" dc:"The unique id of refund" v:"required"`
}
type RefundCancelRes struct {
}
