package model

type Discount struct {
	DiscountTag           string `json:"discountTag,omitempty"`
	DiscountName          string `json:"discountName,omitempty"`
	SavingsAmount         Amount `json:"savingsAmount,omitempty"`
	EstimateSavingsAmount Amount `json:"estimateSavingsAmount,omitempty"`
}
