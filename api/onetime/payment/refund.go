package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentRefundReq struct {
	g.Meta           `path:"/refund/new/{PaymentId}" tags:"OneTime-Payment-Controller" method:"post" summary:"New Payment Refund"`
	PaymentId        string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantId       uint64    `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantRefundId string    `p:"merchantRefundId" dc:"MerchantRefundId" v:"required"`
	Reason           string    `p:"reason" dc:"Reason"`
	Amount           *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type NewPaymentRefundRes struct {
	Status           string `p:"status" dc:"Status"`
	RefundId         string `p:"refundId" dc:"RefundId"`
	MerchantRefundId string `p:"merchantRefundId" dc:"MerchantRefundId"`
	PaymentId        string `p:"paymentId" dc:"PaymentId"`
}
