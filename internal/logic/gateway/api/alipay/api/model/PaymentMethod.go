package model

type PaymentMethod struct {
	PaymentMethodType           string         `json:"paymentMethodType,omitempty"`
	PaymentMethodId             string         `json:"paymentMethodId,omitempty"`
	PaymentMethodMetaData       map[string]any `json:"paymentMethodMetaData,omitempty"`
	CustomerId                  string         `json:"customerId,omitempty"`
	ExtendInfo                  string         `json:"extendInfo,omitempty"`
	RequireIssuerAuthentication bool           `json:"requireIssuerAuthentication,omitempty"`
}
