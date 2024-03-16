package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"Profile" method:"get" summary:"Get Merchant Info"`
}

type GetRes struct {
	Merchant *bean.MerchantSimplify  `json:"merchant" dc:"Merchant"`
	Env      string                  `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd   bool                    `json:"isProd" description:"Check System Env Is Prod, true|false" `
	TimeZone []string                `json:"TimeZone" description:"TimeZone List" `
	Currency []*bean.Currency        `json:"Currency" description:"Currency List" `
	Gateway  []*bean.GatewaySimplify `json:"gateway" description:"Gateway List" `
}

type UpdateReq struct {
	g.Meta      `path:"/update" tags:"Profile" method:"post" summary:"Update Merchant Info"`
	CompanyName string `json:"companyName" description:"company_name"`
	Email       string `json:"email"       description:"email"`
	Address     string `json:"address"     description:"address"`
	CompanyLogo string `json:"companyLogo" description:"company_logo"`
	Phone       string `json:"phone"       description:"phone"`
	TimeZone    string `json:"timeZone" description:"User TimeZone"`
	Host        string `json:"host" description:"User Portal Host"`
}

type UpdateRes struct {
	Merchant *bean.MerchantSimplify `json:"merchant" dc:"Merchant"`
}

type CountryConfigListReq struct {
	g.Meta `path:"/country_config_list" tags:"Profile" method:"post" summary:"Merchant Edit Country Config"`
}
type CountryConfigListRes struct {
	Configs []*bean.MerchantCountryConfigSimplify `json:"configs" description:"Configs"`
}

type EditCountryConfigReq struct {
	g.Meta      `path:"/edit_country_config" tags:"Profile" method:"post" summary:"Merchant Country Config List"`
	CountryCode string `json:"countryCode"  dc:"CountryCode" v:"required"`
	Name        string `json:"name"  dc:"name" `
}
type EditCountryConfigRes struct {
}
