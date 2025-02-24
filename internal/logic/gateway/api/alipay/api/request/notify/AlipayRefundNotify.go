package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayRefundNotify struct {
	AlipayNotify
	RefundStatus          string        `json:"refundStatus,omitempty"`
	RefundRequestId       string        `json:"refundRequestId,omitempty"`
	RefundId              string        `json:"refundId,omitempty"`
	RefundAmount          *model.Amount `json:"refundAmount,omitempty"`
	RefundTime            string        `json:"refundTime,omitempty"`
	GrossSettlementAmount *model.Amount `json:"grossSettlementAmount,omitempty"`
	SettlementQuote       *model.Quote  `json:"settlementQuote,omitempty"`
}
