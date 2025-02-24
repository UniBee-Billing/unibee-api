package model

type AccountBalance struct {
	AccountNo        string  `json:"accountNo,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	AvailableBalance *Amount `json:"availableBalance,omitempty"`
	FrozenBalance    *Amount `json:"frozenBalance,omitempty"`
	TotalBalance     *Amount `json:"totalBalance,omitempty"`
}
