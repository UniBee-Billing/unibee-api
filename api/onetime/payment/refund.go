package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentRefundReq struct {
	g.Meta           `path:"/refund/new/{PaymentId}" tags:"OneTime-Payment" method:"post" summary:"New Payment Refund"`
	PaymentId        string    `in:"path" dc:"PaymentId" v:"required"`
	ExternalRefundId string    `json:"externalRefundId" dc:"ExternalRefundId" v:"required"`
	Amount           *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
	Reason           string    `json:"reason" dc:"Reason"`
}
type NewPaymentRefundRes struct {
	Status           string `json:"status" dc:"Status"`
	RefundId         string `json:"refundId" dc:"RefundId"`
	MerchantRefundId string `json:"merchantRefundId" dc:"ExternalRefundId"`
	PaymentId        string `json:"paymentId" dc:"PaymentId"`
}
