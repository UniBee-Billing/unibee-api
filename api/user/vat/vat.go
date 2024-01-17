package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/vat_gateway"
)

type CountryVatListReq struct {
	g.Meta     `path:"/vat_country_list" tags:"User-Vat-Controller" method:"post" summary:"Vat Country List"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
}
type CountryVatListRes struct {
	VatCountryList []*vat_gateway.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}
