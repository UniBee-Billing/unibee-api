package model

type CustomsInfo struct {
	CustomsCode string `json:"customsCode,omitempty"`
	Region      string `json:"region,omitempty"`
}

type MerchantCustomsInfo struct {
	MerchantCustomsCode string `json:"merchantCustomsCode,omitempty"`
	MerchantCustomsName string `json:"merchantCustomsName,omitempty"`
}
