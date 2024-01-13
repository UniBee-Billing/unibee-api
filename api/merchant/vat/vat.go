package vat

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/logic/vat_gateway"
)

type SetupVatGatewayReq struct {
	g.Meta     `path:"/vat_gateway_setup" tags:"Merchant-Setting-Controller" method:"post" summary:"Vat Gateway 设置"`
	MerchantId int64  `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
	VatName    string `p:"vatName" dc:"vatName, em. vatsense" v:"required|length:4,30#请输入VatName"`
	VatData    string `p:"vatData" dc:"VatData" v:"required|length:4,30#请输入VatData"`
	IsDefault  bool   `p:"IsDefault" d:"false" dc:"IsDefault, default is false" `
}
type SetupVatGatewayRes struct {
}

type CountryVatListReq struct {
	g.Meta     `path:"/vat_country_list" tags:"Merchant-Vat-Controller" method:"post" summary:"Vat Country List"`
	MerchantId int64 `p:"merchantId" dc:"MerchantId" v:"required|length:4,30#请输入商户号"`
}
type CountryVatListRes struct {
	VatCountryList []*vat_gateway.VatCountryRate `json:"vatCountryList" dc:"VatCountryList"`
}
