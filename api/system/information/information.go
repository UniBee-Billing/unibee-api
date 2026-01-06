package information

import (
	"unibee/api/bean"
	"unibee/api/bean/detail"

	"github.com/gogf/gf/v2/frame/g"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"System-Information" method:"get" summary:"Get System Information"`
}

type GetRes struct {
	Env                string            `json:"env" description:"System Env, em: daily|stage|local|prod" `
	Mode               string            `json:"mode" description:"System Mode" `
	BuildVersion       string            `json:"buildVersion" description:"System Build Version" `
	IsProd             bool              `json:"isProd" description:"Check System Env Is Prod, true|false" `
	SupportTimeZone    []string          `json:"supportTimeZone" description:"Support TimeZone List" `
	SupportCurrency    []*bean.Currency  `json:"supportCurrency" description:"Support Currency List" `
	SupportLanguage    []*Language       `json:"supportLanguage" description:"Support Language List" `
	SupportCountryCode []string          `json:"supportCountryCode" description:"Support Country Code List (ISO 3166-1 alpha-2)" `
	Gateway            []*detail.Gateway `json:"gateway" description:"Support Gateway List" `
	OAuth              *OAuthConfig      `json:"oauth" description:"OAuth Configuration" `
}

type Language struct {
	Code     string `json:"code" description:"Language code (e.g., en, ru, vi, cn, pt)"`
	FullName string `json:"fullName" description:"Language full name (e.g., English, Russian, Vietnamese, Chinese, Portuguese)"`
}

type OAuthConfig struct {
	TokenSecret        string `json:"tokenSecret" description:"OAuth token secret"`
	GoogleClientId     string `json:"googleClientId" description:"Google OAuth client ID"`
	GoogleClientSecret string `json:"googleClientSecret" description:"Google OAuth client secret"`
	GithubClientId     string `json:"githubClientId" description:"GitHub OAuth client ID"`
	//GithubClientSecret string `json:"githubClientSecret" description:"GitHub OAuth client secret"`
}

type SendMockMQReq struct {
	g.Meta  `path:"/send_mock_mq" tags:"System-Information" method:"get" summary:"Send Mock MQ Message"`
	Message string `json:"message" description:"Send Mock MQ Message"`
}

type SendMockMQRes struct {
}
