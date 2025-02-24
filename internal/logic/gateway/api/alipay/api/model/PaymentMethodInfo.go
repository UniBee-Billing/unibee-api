package model

type PaymentMethodInfo struct {
	PaymentMethodType   string `json:"paymentMethodType,omitempty"`
	PaymentMethodDetail string `json:"paymentMethodDetail,omitempty"`
	Enabled             bool   `json:"enabled,omitempty"`
	Preferred           bool   `json:"preferred,omitempty"`
	ExtendInfo          string `json:"extendInfo,omitempty"`
}
