package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"Merchant" method:"get" summary:"Get Profile"`
}

type GetRes struct {
	Merchant             *bean.Merchant               `json:"merchant" dc:"Merchant"`
	MerchantMember       *detail.MerchantMemberDetail `json:"merchantMember" dc:"MerchantMember"`
	Env                  string                       `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd               bool                         `json:"isProd" description:"Check System Env Is Prod, true|false" `
	TimeZone             []string                     `json:"TimeZone" description:"TimeZone List" `
	Currency             []*bean.Currency             `json:"Currency" description:"Currency List" `
	Gateways             []*detail.Gateway            `json:"gateways" description:"Gateway List" `
	ExchangeRateApiKey   string                       `json:"exchangeRateApiKey" description:"ExchangeRateApiKey" `
	OpenApiKey           string                       `json:"openApiKey" description:"OpenApiKey" `
	SendGridKey          string                       `json:"sendGridKey" description:"SendGridKey" `
	VatSenseKey          string                       `json:"vatSenseKey" description:"VatSenseKey" `
	EmailSender          *bean.Sender                 `json:"emailSender" description:"EmailSender" `
	SegmentServerSideKey string                       `json:"segmentServerSideKey" description:"SegmentServerSideKey" `
	SegmentUserPortalKey string                       `json:"segmentUserPortalKey" description:"SegmentUserPortalKey" `
	IsOwner              bool                         `json:"isOwner" description:"Check Member is Owner" `
	MemberRoles          []*bean.MerchantRole         `json:"MemberRoles" description:"The member role list'" `
}

type UpdateReq struct {
	g.Meta      `path:"/update" tags:"Merchant" method:"post" summary:"Update Profile"`
	CompanyName string `json:"companyName" description:"company_name"`
	Email       string `json:"email"       description:"email"`
	Address     string `json:"address"     description:"address"`
	CompanyLogo string `json:"companyLogo" description:"company_logo"`
	Phone       string `json:"phone"       description:"phone"`
	TimeZone    string `json:"timeZone" description:"User TimeZone"`
	Host        string `json:"host" description:"User Portal Host"`
}

type UpdateRes struct {
	Merchant *bean.Merchant `json:"merchant" dc:"Merchant"`
}

type CountryConfigListReq struct {
	g.Meta `path:"/country_config_list" tags:"Merchant" method:"post" summary:"Edit Country Config"`
}
type CountryConfigListRes struct {
	Configs []*bean.MerchantCountryConfig `json:"configs" description:"Configs"`
}

type EditCountryConfigReq struct {
	g.Meta      `path:"/edit_country_config" tags:"Merchant" method:"post" summary:"Get Country Config List"`
	CountryCode string `json:"countryCode"  dc:"CountryCode" v:"required"`
	Name        string `json:"name"  dc:"name" `
	VatEnable   *bool  `json:"vatEnable"  dc:"VatEnable, Default true" `
}
type EditCountryConfigRes struct {
}

type NewApiKeyReq struct {
	g.Meta `path:"/new_apikey" tags:"Merchant" method:"post" summary:"Generate New APIKey" dc:"Generate new apikey, The old one expired in one hour"`
}
type NewApiKeyRes struct {
	ApiKey string `json:"apiKey" description:"ApiKey"`
}
