package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type SetupVatGatewayReq struct {
	g.Meta      `path:"/vat_gateway_setup" tags:"Merchant-Setting-Controller" method:"post" summary:"Vat Gateway Settings"`
	GatewayName string `p:"gatewayName" dc:"GatewayName, em. vatsense" v:"required"`
	Data        string `p:"data" dc:"Data" v:"required"`
	IsDefault   bool   `p:"IsDefault" d:"true" dc:"IsDefault, default is true" `
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
