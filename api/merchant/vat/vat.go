package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type SetupVatGatewayReq struct {
	g.Meta    `path:"/vat_gateway_setup" tags:"Merchant-Setting-Controller" method:"post" summary:"Vat Gateway Settings"`
	VatName   string `p:"vatName" dc:"vatName, em. vatsense" v:"required"`
	VatData   string `p:"vatData" dc:"VatData" v:"required"`
	IsDefault bool   `p:"IsDefault" d:"false" dc:"IsDefault, default is false" `
}
type SetupVatGatewayRes struct {
}

type CountryVatListReq struct {
	g.Meta     `path:"/vat_country_list" tags:"Merchant-Vat-Controller" method:"post" summary:"Vat Country List"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
}
type CountryVatListRes struct {
	VatCountryList []*ro.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}
