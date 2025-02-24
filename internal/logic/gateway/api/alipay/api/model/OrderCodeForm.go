package model

type OrderCodeForm struct {
	PaymentMethodType string       `json:"paymentMethodType,omitempty"`
	ExpireTime        string       `json:"expireTime,omitempty"`
	CodeDetails       []CodeDetail `json:"codeDetails,omitempty"`
	ExtendInfo        string       `json:"extendInfo,omitempty"`
}
