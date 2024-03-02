package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type SetupGatewayReq struct {
	g.Meta      `path:"/setup_gateway" tags:"Merchant-Vat-Controller" method:"post" summary:"Vat Gateway Setup"`
	GatewayName string `p:"gatewayName" dc:"GatewayName, em. vatsense" v:"required"`
	Data        string `p:"data" dc:"Data" v:"required"`
	IsDefault   bool   `p:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}
type SetupGatewayRes struct {
}

type CountryListReq struct {
	g.Meta     `path:"/country_list" tags:"Merchant-Vat-Controller" method:"post" summary:"Vat Country List"`
	MerchantId uint64 `p:"merchantId" dc:"MerchantId" v:"required"`
}
type CountryListRes struct {
	VatCountryList []*ro.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}
