package profile

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/logic/middleware/license"
)

type GetLicenseReq struct {
	g.Meta `path:"/get_license" tags:"Merchant" method:"get" summary:"Get License"`
}

type GetLicenseRes struct {
	Merchant           *bean.Merchant   `json:"merchant" dc:"Merchant"`
	License            *license.License `json:"license" description:"The merchant license" `
	APIRateLimit       int              `json:"APIRateLimit" dc:"APIRateLimit"`
	MemberLimit        int              `json:"memberLimit" dc:"MemberLimit, -1=Unlimited"`
	CurrentMemberCount int              `json:"currentMemberCount" dc:"CurrentMemberCount"`
}

type GetLicenseUpdateUrlReq struct {
	g.Meta    `path:"/get_license_update_url" tags:"Merchant" method:"get,post" summary:"Get License Update Url"`
	PlanId    int64  `json:"planId" dc:"Id of plan to update" dc:"Id of plan to update"`
	ReturnUrl string `json:"returnUrl"  dc:"ReturnUrl"`
	CancelUrl string `json:"cancelUrl" dc:"CancelUrl"`
}

type GetLicenseUpdateUrlRes struct {
	Url string `json:"url" dc:"Url"`
}

type GetReq struct {
	g.Meta `path:"/get" tags:"Merchant" method:"get" summary:"Get Profile"`
}

type GetRes struct {
	Merchant                     *bean.Merchant                      `json:"merchant" dc:"Merchant"`
	MerchantMember               *detail.MerchantMemberDetail        `json:"merchantMember" dc:"MerchantMember"`
	Env                          string                              `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd                       bool                                `json:"isProd" description:"Check System Env Is Prod, true|false" `
	TimeZone                     []string                            `json:"TimeZone" description:"TimeZone List" `
	DefaultCurrency              string                              `json:"defaultCurrency" description:"Default Currency" `
	Currency                     []*bean.Currency                    `json:"Currency" description:"Currency List" `
	Gateways                     []*detail.Gateway                   `json:"gateways" description:"Gateway List" `
	ExchangeRateApiKey           string                              `json:"exchangeRateApiKey" description:"ExchangeRateApiKey" `
	OpenAPIHost                  string                              `json:"openApiHost" description:"OpenApi Host"`
	OpenAPIKey                   string                              `json:"openApiKey" description:"OpenAPIKey" `
	SendGridKey                  string                              `json:"sendGridKey" description:"SendGridKey" `
	VatSenseKey                  string                              `json:"vatSenseKey" description:"VatSenseKey" `
	EmailSender                  *bean.Sender                        `json:"emailSender" description:"EmailSender" `
	SegmentServerSideKey         string                              `json:"segmentServerSideKey" description:"SegmentServerSideKey" `
	SegmentUserPortalKey         string                              `json:"segmentUserPortalKey" description:"SegmentUserPortalKey" `
	GlobalTOPTEnabled            bool                                `json:"globalTOPTEnabled" description:"GlobalTOPTEnabled" `
	QuickBooksCompanyName        string                              `json:"quickBooksCompanyName" description:"QuickBooksCompanyName" `
	QuickBooksLastSynchronized   string                              `json:"quickBooksLastSynchronized" description:"QuickBooksLastSynchronized" `
	QuickBooksLastSyncError      string                              `json:"quickBooksLastSyncError" description:"QuickBooksLastSyncError" `
	IsOwner                      bool                                `json:"isOwner" description:"Check Member is Owner" `
	MemberRoles                  []*bean.MerchantRole                `json:"MemberRoles" description:"The member role list'" `
	CloudFeatureAnalyticsEnabled bool                                `json:"cloudFeatureAnalyticsEnabled" description:"Analytics Feature Enabled For Cloud Version"`
	AnalyticsHost                string                              `json:"analyticsHost" description:"Analytics Host"`
	MultiCurrencies              []*bean.MerchantMultiCurrencyConfig `json:"multiCurrencyConfigs"  dc:"Merchant's MultiCurrency Configs" `
}

type UpdateReq struct {
	g.Meta              `path:"/update" tags:"Merchant" method:"post" summary:"Update Profile"`
	CompanyName         string `json:"companyName" description:"company_name"`
	Email               string `json:"email"       description:"email"`
	Address             string `json:"address"     description:"address"`
	CompanyLogo         string `json:"companyLogo" description:"company_logo"`
	Phone               string `json:"phone"       description:"phone"`
	TimeZone            string `json:"timeZone" description:"User TimeZone"`
	Host                string `json:"host" description:"User Portal Host"`
	CountryCode         string `json:"countryCode" dc:"Country Code"`
	CountryName         string `json:"countryName" dc:"Country Name"`
	CompanyVatNumber    string `json:"companyVatNumber" dc:"Country Vat Number"`
	CompanyRegistryCode string `json:"companyRegistryCode" dc:"Country Registry Code"`
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

type EditTotpConfigReq struct {
	g.Meta   `path:"/edit_totp_config" tags:"Merchant" method:"post" summary:"Admin Edit 2FA Config"`
	Activate bool `json:"activate"   description:"activate 2FA for all members, all members need 2FA while login if activate, otherwise not"`
}

type EditTotpConfigRes struct {
}

type NewApiKeyReq struct {
	g.Meta `path:"/new_apikey" tags:"Merchant" method:"post" summary:"Generate New APIKey" dc:"Generate new apikey, The old one expired in one hour"`
}
type NewApiKeyRes struct {
	ApiKey string `json:"apiKey" description:"ApiKey"`
}

type SetupMultiCurrenciesReq struct {
	g.Meta          `path:"/setup_multi_currencies" tags:"Merchant" method:"post" summary:"Multi Currencies Setup"`
	MultiCurrencies []*bean.MerchantMultiCurrencyConfig `json:"multiCurrencyConfigs"  dc:"Merchant's MultiCurrencies" `
}
type SetupMultiCurrenciesRes struct {
	MultiCurrencies []*bean.MerchantMultiCurrencyConfig `json:"multiCurrencyConfigs"  dc:"Merchant's MultiCurrencies" `
}

type AmountMultiCurrenciesExchangeReq struct {
	g.Meta   `path:"/amount_multi_currencies_exchange" tags:"Merchant" method:"post" summary:"Amount Multi Currencies Exchange"`
	Amount   int64  `json:"amount"   dc:"Amount"   v:"required" `
	Currency string `json:"currency"   dc:"The Default Currency" v:"required" `
}
type AmountMultiCurrenciesExchangeRes struct {
	MultiCurrencies []*bean.PlanMultiCurrency `json:"multiCurrencyConfigs"  dc:"Merchant's MultiCurrencies" `
}
