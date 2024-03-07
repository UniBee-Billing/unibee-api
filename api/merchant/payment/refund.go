package payment

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NewPaymentRefundReq struct {
	g.Meta           `path:"/refund/new" tags:"Payment" method:"post" summary:"New Payment Refund"`
	PaymentId        string `json:"paymentId" dc:"PaymentId" v:"required"`
	ExternalRefundId string `json:"externalRefundId" dc:"ExternalRefundId" v:"required"`
	RefundAmount     int64  `json:"refundAmount"   in:"query" dc:"RefundAmount, Cent" v:"required"`
	Currency         string `json:"currency"   in:"query" dc:"Currency"  v:"required"`
	Reason           string `json:"reason" dc:"Reason"`
}
type NewPaymentRefundRes struct {
	Status           int    `json:"status" dc:"Status,10-create|20-success|30-Failed|40-Reverse"`
	RefundId         string `json:"refundId" dc:"RefundId"`
	ExternalRefundId string `json:"externalRefundId" dc:"ExternalRefundId"`
	PaymentId        string `json:"paymentId" dc:"PaymentId"`
}

type RefundDetailReq struct {
	g.Meta   `path:"/refund/detail" tags:"Payment" method:"get" summary:"Query Payment Refund Detail"`
	RefundId string `json:"refundId" dc:"RefundId"`
}

type RefundDetailRes struct {
}

type RefundListReq struct {
	g.Meta         `path:"/refund/list" tags:"Payment" method:"get" summary:"Query Payment Refund List"`
	PaymentId      string `json:"paymentId" dc:"PaymentId" v:"required"`
	Status         int    `json:"status" dc:"Status,10-create|20-success|30-Failed|40-Reverse"`
	GatewayId      uint64 `json:"gatewayId"   dc:"GatewayId"`
	ExternalUserId string `json:"externalUserId" dc:"ExternalUserId"`
	Email          string `json:"email" dc:"Email"`
	Currency       string `json:"currency" dc:"Currency"`
}
type RefundListRes struct {
}
