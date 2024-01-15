package base

import entity "go-oversea-pay/internal/model/entity/oversea_pay"

type ValidResult struct {
	Valid          bool
	VatNumber      string
	CountryCode    string
	CompanyName    string
	CompanyAddress string
}

type Gateway interface {
	GetVatName() string
	ListAllCountries() ([]*entity.CountryRate, error)
	ListAllRates() ([]*entity.CountryRate, error)
	ValidateVatNumber(vatNumber string, requesterVatNumber string) (*ValidResult, error)
	ValidateEoriNumber(number string) (*ValidResult, error)
}
