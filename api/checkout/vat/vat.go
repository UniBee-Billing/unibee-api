package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type CountryListReq struct {
	g.Meta     `path:"/country_list" tags:"Checkout" method:"get,post" summary:"Vat Country List"`
	MerchantId uint64 `json:"merchantId" description:""  v:"required"`
}
type CountryListRes struct {
	VatCountryList []*bean.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta     `path:"/vat_number_validate" tags:"Checkout" method:"post" summary:"Vat Number Validate"`
	MerchantId uint64 `json:"merchantId" description:""  v:"required"`
	VatNumber  string `json:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *bean.ValidResult `json:"vatNumberValidate"`
}
