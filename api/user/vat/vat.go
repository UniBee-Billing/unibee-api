package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/vat_gateway"
	"go-oversea-pay/internal/logic/vat_gateway/base"
)

type CountryVatListReq struct {
	g.Meta     `path:"/vat_country_list" tags:"User-Vat-Controller" method:"post" summary:"Vat Country List"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required"`
}
type CountryVatListRes struct {
	VatCountryList []*vat_gateway.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta     `path:"/vat_number_validate" tags:"User-Vat-Controller" method:"post" summary:"Vat Number Validate"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required"`
	VatNumber  string `p:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *base.ValidResult `json:"vatNumberValidate"`
}
