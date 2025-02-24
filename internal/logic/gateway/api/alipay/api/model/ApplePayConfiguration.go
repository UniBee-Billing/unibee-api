package model

type ApplePayConfiguration struct {
	RequiredBillingContactFields  []string `json:"requiredBillingContactFields,omitempty"`
	RequiredShippingContactFields []string `json:"requiredShippingContactFields,omitempty"`
	ButtonsBundled                bool     `json:"buttonsBundled,omitempty"`
	ApplePayToken                 string   `json:"applePayToken,omitempty"`
}
