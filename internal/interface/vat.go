package _interface

import (
	"unibee/internal/logic/gateway/ro"
	entity "unibee/internal/model/entity/oversea_pay"
)

type VATGateway interface {
	GetGatewayName() string
	ListAllCountries() ([]*entity.CountryRate, error)
	ListAllRates() ([]*entity.CountryRate, error)
	ValidateVatNumber(vatNumber string, requesterVatNumber string) (*ro.ValidResult, error)
	ValidateEoriNumber(number string) (*ro.ValidResult, error)
}
