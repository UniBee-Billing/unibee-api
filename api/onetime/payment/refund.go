package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentRefundReq struct {
	g.Meta           `path:"/refund/new/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"New Payment Refund"`
	PaymentId        string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantRefundId string    `json:"merchantRefundId" dc:"MerchantRefundId" v:"required"`
	Reason           string    `json:"reason" dc:"Reason"`
	Amount           *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type NewPaymentRefundRes struct {
	Status           string `json:"status" dc:"Status"`
	RefundId         string `json:"refundId" dc:"RefundId"`
	MerchantRefundId string `json:"merchantRefundId" dc:"MerchantRefundId"`
	PaymentId        string `json:"paymentId" dc:"PaymentId"`
}
