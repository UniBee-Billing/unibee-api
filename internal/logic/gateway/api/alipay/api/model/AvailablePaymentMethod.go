package model

type AvailablePaymentMethod struct {
	PaymentMethodTypeList []PaymentMethodTypeItem `json:"paymentMethodTypeList,omitempty"`
}

type PaymentMethodTypeItem struct {
	PaymentMethodType  string `json:"paymentMethodType,omitempty"`
	PaymentMethodOrder int    `json:"paymentMethodOrder,omitempty"`
	ExpressCheckout    bool   `json:"expressCheckout,omitempty"`
}
