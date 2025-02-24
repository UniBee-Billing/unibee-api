package model

type CardInfo struct {
	CardNo              string        `json:"cardNo,omitempty"`
	CardBrand           CardBrand     `json:"cardBrand,omitempty"`
	CardToken           string        `json:"cardToken,omitempty"`
	IssuingCountry      string        `json:"issuingCountry,omitempty"`
	Funding             string        `json:"funding,omitempty"`
	PaymentMethodRegion string        `json:"paymentMethodRegion,omitempty"`
	ThreeDSResult       ThreeDSResult `json:"threeDSResult,omitempty"`
}
