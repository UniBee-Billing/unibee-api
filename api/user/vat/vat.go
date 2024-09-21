package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type CountryListReq struct {
	g.Meta `path:"/country_list" tags:"User-Vat" method:"get,post" summary:"Vat Country List"`
}
type CountryListRes struct {
	VatCountryList []*bean.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta    `path:"/vat_number_validate" tags:"User-Vat" method:"post" summary:"Vat Number Validate"`
	VatNumber string `json:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *bean.ValidResult `json:"vatNumberValidate"`
}
