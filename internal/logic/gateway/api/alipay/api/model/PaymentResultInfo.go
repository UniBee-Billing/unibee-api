package model

type PaymentResultInfo struct {
	CardNo               string        `json:"cardNo,omitempty"`
	CardBrand            string        `json:"cardBrand,omitempty"`
	CardToken            string        `json:"cardToken,omitempty"`
	IssuingCountry       string        `json:"issuingCountry,omitempty"`
	Funding              string        `json:"funding,omitempty"`
	PaymentMethodRegion  string        `json:"paymentMethodRegion,omitempty"`
	ThreeDSResult        ThreeDSResult `json:"threeDSResult,omitempty"`
	AvsResultRaw         string        `json:"avsResultRaw,omitempty"`
	CvvResultRaw         string        `json:"cvvResultRaw,omitempty"`
	NetworkTransactionId string        `json:"networkTransactionId,omitempty"`
	CreditPayPlan        CreditPayPlan `json:"creditPayPlan,omitempty"`
}

type ThreeDSResult struct {
	ThreeDSVersion  string `json:"threeDSVersion,omitempty"`
	Eci             string `json:"eci,omitempty"`
	Cavv            string `json:"cavv,omitempty"`
	DsTransactionId string `json:"dsTransactionId,omitempty"`
	Xid             string `json:"xid,omitempty"`
}
