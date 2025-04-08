package information

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	"unibee/api/bean/detail"
)

type GetReq struct {
	g.Meta `path:"/get" tags:"System-Information" method:"get" summary:"Get System Information"`
}

type GetRes struct {
	Env             string            `json:"env" description:"System Env, em: daily|stage|local|prod" `
	Mode            string            `json:"mode" description:"System Mode" `
	BuildVersion    string            `json:"buildVersion" description:"System Build Version" `
	IsProd          bool              `json:"isProd" description:"Check System Env Is Prod, true|false" `
	SupportTimeZone []string          `json:"supportTimeZone" description:"Support TimeZone List" `
	SupportCurrency []*bean.Currency  `json:"supportCurrency" description:"Support Currency List" `
	Gateway         []*detail.Gateway `json:"gateway" description:"Support Currency List" `
}

type SendMockMQReq struct {
	g.Meta  `path:"/send_mock_mq" tags:"System-Information" method:"get" summary:"Send Mock MQ Message"`
	Message string `json:"message" description:"Send Mock MQ Message"`
}

type SendMockMQRes struct {
}
