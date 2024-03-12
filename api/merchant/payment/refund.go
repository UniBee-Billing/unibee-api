package payment

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type NewPaymentRefundReq struct {
	g.Meta           `path:"/refund/new" tags:"Payment" method:"post" summary:"New Payment Refund"`
	PaymentId        string            `json:"paymentId" dc:"PaymentId" v:"required"`
	ExternalRefundId string            `json:"externalRefundId" dc:"ExternalRefundId" v:"required"`
	RefundAmount     int64             `json:"refundAmount" dc:"RefundAmount, Cent" v:"required"`
	Currency         string            `json:"currency" dc:"Currency" v:"required"`
	Reason           string            `json:"reason" dc:"Reason"`
	Metadata         map[string]string `json:"metadata" dc:"Metadata，Map"`
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
	RefundDetail *RefundDetail `json:"refundDetail" dc:"RefundDetail"`
}

type RefundDetail struct {
	User    *bean.UserAccountSimplify `json:"user" dc:"user"`
	Payment *bean.PaymentSimplify     `json:"payment" dc:"Payment"`
	Refund  *bean.RefundSimplify      `json:"refund" dc:"Refund"`
}

type RefundListReq struct {
	g.Meta    `path:"/refund/list" tags:"Payment" method:"get" summary:"Query Payment Refund List"`
	PaymentId string `json:"paymentId" dc:"PaymentId" v:"required"`
	Status    int    `json:"status" dc:"Status,10-create|20-success|30-Failed|40-Reverse"`
	GatewayId uint64 `json:"gatewayId" dc:"GatewayId"`
	UserId    int64  `json:"userId" dc:"UserId"`
	Email     string `json:"email" dc:"Email"`
	Currency  string `json:"currency" dc:"Currency"`
}
type RefundListRes struct {
	RefundDetails []*RefundDetail `json:"refundDetails" dc:"RefundDetails"`
}
