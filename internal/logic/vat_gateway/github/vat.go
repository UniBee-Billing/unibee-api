/*
Package vat helps you deal with European VAT in Go.

It offers VAT number validation using the VIES VAT validation API & VAT rates retrieval using jsonvat.com

Validate a VAT number

	validity := vat.ValidateNumber("NL123456789B01")

Get VAT rate that is currently in effect for a given country

	c, _ := vat.GetCountryRates("NL")
	r, _ := c.GetRate("standard")
*/
package vat

import (
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"go-oversea-pay/internal/logic/channel/ro"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"go-oversea-pay/utility"
)

// ErrServiceUnavailable will be returned when VIES VAT validation API or jsonvat.com is unreachable.
var ErrServiceUnavailable = errors.New("vat: service is unreachable")

// ServiceTimeout indicates the number of seconds before a service request times out.
var ServiceTimeout = 10

type Github struct {
	Password string
	Name     string
}

func (g Github) GetGatewayName() string {
	return g.Name
}

func (g Github) ListAllCountries() ([]*entity.CountryRate, error) {
	return []*entity.CountryRate{}, nil
}

func (g Github) ListAllRates() ([]*entity.CountryRate, error) {
	rates, err := FetchRates()
	if err != nil {
		return nil, err
	}
	var list []*entity.CountryRate
	for _, rate := range rates {
		r, _ := rate.GetRate("standard")
		var vat = 0
		if r > 0 {
			vat = 1
		} else {
			vat = 2
		}
		list = append(list, &entity.CountryRate{
			Gateway:               g.GetGatewayName(),
			CountryCode:           rate.CountryCode,
			CountryName:           "", // todo mark
			StandardTaxPercentage: int64(r * 100),
			Vat:                   vat,
			Other:                 utility.FormatToJsonString(rate),
		})
	}
	return list, nil
}

func (g Github) ValidateVatNumber(vatNumber string, requesterVatNumber string) (*ro.ValidResult, error) {
	format, err := ValidateNumberFormat(vatNumber)
	if err != nil {
		return nil, err
	}
	if !format {
		return nil, gerror.New(vatNumber + " is not valid format")
	}
	valid, err := ValidateNumberExistenceV2(vatNumber)
	if err != nil {
		return nil, err
	}
	if !valid.Valid {
		return nil, gerror.New(vatNumber + " is not valid")
	}
	return &ro.ValidResult{
		Valid:          valid.Valid,
		VatNumber:      valid.VATNumber,
		CountryCode:    valid.CountryCode,
		CompanyName:    valid.Name,
		CompanyAddress: valid.Address,
	}, nil
}

func (g Github) ValidateEoriNumber(number string) (*ro.ValidResult, error) {
	return nil, gerror.New("not support")
}
