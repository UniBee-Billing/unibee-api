package model

type RefundDetail struct {
	RefundAmount Amount         `json:"refundAmount,omitempty"`
	RefundFrom   RefundFromType `json:"refundFrom,omitempty"`
}

type RefundFromType string

const (
	//SELLER,
	//MARKETPLACE,
	//UNSETTLED_FUNDS;
	RefundFromTypeSELLER         RefundFromType = "SELLER"
	RefundFromTypeMARKETPLACE    RefundFromType = "MARKETPLACE"
	RefundFromTypeUNSETTLEDFUNDS RefundFromType = "UNSETTLED_FUNDS"
)
