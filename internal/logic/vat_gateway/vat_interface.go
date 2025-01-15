package vat_gateway

import (
	"unibee/api/bean"
	entity "unibee/internal/model/entity/default"
)

type VATGateway interface {
	GetGatewayName() string
	ListAllCountries() ([]*entity.CountryRate, error)
	ListAllRates() ([]*entity.CountryRate, error)
	ValidateVatNumber(vatNumber string, requesterVatNumber string) (*bean.ValidResult, error)
	ValidateEoriNumber(number string) (*bean.ValidResult, error)
}
