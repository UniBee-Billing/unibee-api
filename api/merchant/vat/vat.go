package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type SetupGatewayReq struct {
	g.Meta      `path:"/setup_gateway" tags:"Vat" method:"post" summary:"Vat Gateway Setup"`
	GatewayName string `json:"gatewayName" dc:"GatewayName, em. vatsense" v:"required"`
	Data        string `json:"data" dc:"Data" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}
type SetupGatewayRes struct {
}

type CountryListReq struct {
	g.Meta     `path:"/country_list" tags:"Vat" method:"get,post" summary:"Vat Country List"`
	MerchantId uint64 `json:"merchantId" dc:"MerchantId" v:"required"`
}
type CountryListRes struct {
	VatCountryList []*bean.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}
