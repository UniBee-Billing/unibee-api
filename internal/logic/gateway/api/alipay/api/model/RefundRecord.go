package model

type RefundRecord struct {
	ReferenceOrderId string  `json:"referenceOrderId,omitempty"`
	ReferenceGoodsId string  `json:"referenceGoodsId,omitempty"`
	Amount           *Amount `json:"amount,omitempty"`
	RefundReason     string  `json:"refundReason,omitempty"`
	RefundTime       string  `json:"refundTime,omitempty"`
}
