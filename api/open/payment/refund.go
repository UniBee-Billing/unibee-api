package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type RefundsReq struct {
	g.Meta           `path:"/refunds/{PaymentId}" tags:"Open-Payment-Controller" method:"post" summary:"Refund Create"`
	PaymentId        string    `in:"path" dc:"PaymentId" v:"required"`
	MerchantId       int64     `p:"merchantId" dc:"MerchantId" v:"required"`
	MerchantRefundId string    `p:"merchantRefundId" dc:"MerchantRefundId" v:"required"`
	Reason           string    `p:"reason" dc:"Reason"`
	Amount           *AmountVo `json:"amount"   in:"query" dc:"Amount, Cent" v:"required"`
}
type RefundsRes struct {
	Status           string `p:"status" dc:"Status"`
	RefundId         string `p:"refundId" dc:"RefundId"`
	MerchantRefundId string `p:"merchantRefundId" dc:"MerchantRefundId"`
	PaymentId        string `p:"paymentId" dc:"PaymentId"`
}
