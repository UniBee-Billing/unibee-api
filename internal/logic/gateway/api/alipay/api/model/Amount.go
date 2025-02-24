package model

type Amount struct {
	Value    string `json:"value,omitempty"`
	Currency string `json:"currency,omitempty"`
}

func NewAmount(value string, currency string) *Amount {
	return &Amount{
		Value:    value,
		Currency: currency,
	}
}
