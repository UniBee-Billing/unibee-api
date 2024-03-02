package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type CountryListReq struct {
	g.Meta `path:"/country_list" tags:"User-Vat-Controller" method:"post" summary:"Vat Country List"`
}
type CountryListRes struct {
	VatCountryList []*ro.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta    `path:"/number_validate" tags:"User-Vat-Controller" method:"post" summary:"Vat Number Validate"`
	VatNumber string `p:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *ro.ValidResult `json:"vatNumberValidate"`
}
