package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type SetupGatewayReq struct {
	g.Meta      `path:"/setup_gateway" tags:"Vat Gateway" method:"post" summary:"Vat Gateway Setup"`
	GatewayName string `json:"gatewayName" dc:"GatewayName, em. vatsense" v:"required"`
	Data        string `json:"data" dc:"Data" v:"required"`
	IsDefault   bool   `json:"IsDefault" d:"true" dc:"IsDefault, default is true" `
}
type SetupGatewayRes struct {
	Data string `json:"data" dc:"Data" dc:"The hide star data"`
}

type InitDefaultGatewayReq struct {
	g.Meta `path:"/init_default_gateway" tags:"Vat Gateway" method:"post" summary:"Init Default Vat Gateway"`
}
type InitDefaultGatewayRes struct {
}

type CountryListReq struct {
	g.Meta `path:"/country_list" tags:"Vat Gateway" method:"get,post" summary:"Get Vat Country List"`
}
type CountryListRes struct {
	VatCountryList []*bean.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}

type NumberValidateReq struct {
	g.Meta    `path:"/vat_number_validate" tags:"Vat Gateway" method:"post" summary:"Vat Number Validation"`
	VatNumber string `json:"vatNumber" dc:"VatNumber" v:"required"`
}
type NumberValidateRes struct {
	VatNumberValidate *bean.ValidResult `json:"vatNumberValidate"`
}
