package base

import entity "go-oversea-pay/internal/model/entity/oversea_pay"

type ValidResult struct {
	Valid          bool
	VatNumber      string
	CountryCode    string
	CompanyName    string
	CompanyAddress string
}

const (
	VAT_IMPLEMENT_NAMES = "vatsense"
)

type Gateway interface {
	SetVatName(name string)
	GetVatName() string
	SetVatSignData(data string)
	ListAllCountries() ([]*entity.CountryRate, error)
	ListAllRates() ([]*entity.CountryRate, error)
	ValidateVatNumber(varNumber string, requesterVatNumber string) (*ValidResult, error)
	ValidateEoriNumber(number string) (*ValidResult, error)
}
