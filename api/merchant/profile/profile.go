package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"Profile" method:"get" summary:"MerchantProfile"`
}

type GetRes struct {
	Merchant             *bean.MerchantSimplify       `json:"merchant" dc:"Merchant"`
	MerchantMember       *detail.MerchantMemberDetail `json:"merchantMember" dc:"MerchantMember"`
	Env                  string                       `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd               bool                         `json:"isProd" description:"Check System Env Is Prod, true|false" `
	TimeZone             []string                     `json:"TimeZone" description:"TimeZone List" `
	Currency             []*bean.Currency             `json:"Currency" description:"Currency List" `
	Gateways             []*bean.GatewaySimplify      `json:"gateways" description:"Gateway List" `
	OpenApiKey           string                       `json:"openApiKey" description:"OpenApiKey" `
	SendGridKey          string                       `json:"sendGridKey" description:"SendGridKey" `
	VatSenseKey          string                       `json:"vatSenseKey" description:"VatSenseKey" `
	SegmentServerSideKey string                       `json:"segmentServerSideKey" description:"SegmentServerSideKey" `
	SegmentUserPortalKey string                       `json:"segmentUserPortalKey" description:"SegmentUserPortalKey" `
	IsOwner              bool                         `json:"isOwner" description:"Check Member is Owner" `
	MemberRoles          []*bean.MerchantRoleSimplify `json:"MemberRoles" description:"The member's role list'" `
}

type UpdateReq struct {
	g.Meta      `path:"/update" tags:"Profile" method:"post" summary:"UpdateMerchantProfile"`
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
	g.Meta `path:"/country_config_list" tags:"Profile" method:"post" summary:"EditCountryConfig"`
}
type CountryConfigListRes struct {
	Configs []*bean.MerchantCountryConfigSimplify `json:"configs" description:"Configs"`
}

type EditCountryConfigReq struct {
	g.Meta      `path:"/edit_country_config" tags:"Profile" method:"post" summary:"CountryConfigList"`
	CountryCode string `json:"countryCode"  dc:"CountryCode" v:"required"`
	Name        string `json:"name"  dc:"name" `
	VatEnable   *bool  `json:"vatEnable"  dc:"VatEnable, Default true" `
}
type EditCountryConfigRes struct {
}

type NewApiKeyReq struct {
	g.Meta `path:"/new_apikey" tags:"Profile" method:"post" summary:"GenerateNewApiKey" dc:"generate new apikey, The old one expired in one hour"`
}
type NewApiKeyRes struct {
	ApiKey string `json:"apiKey" description:"ApiKey"`
}
