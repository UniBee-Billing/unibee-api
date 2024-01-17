package base

import entity "go-oversea-pay/internal/model/entity/oversea_pay"

type ValidResult struct {
	Valid           bool   `json:"valid"           `
	VatNumber       string `json:"vatNumber"           `
	CountryCode     string `json:"countryCode"           `
	CompanyName     string `json:"companyName"           `
	CompanyAddress  string `json:"companyAddress"           `
	ValidateMessage string `json:"validateMessage"           `
}

type Gateway interface {
	GetVatName() string
	ListAllCountries() ([]*entity.CountryRate, error)
	ListAllRates() ([]*entity.CountryRate, error)
	ValidateVatNumber(vatNumber string, requesterVatNumber string) (*ValidResult, error)
	ValidateEoriNumber(number string) (*ValidResult, error)
}
