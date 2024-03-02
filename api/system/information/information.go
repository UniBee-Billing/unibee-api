package information

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type GetReq struct {
	g.Meta `path:"/information" tags:"System-Information-Controller" method:"post" summary:"Get System Information"`
}

type GetRes struct {
	Env             string                `json:"env" description:"System Env, em: daily|stage|local|prod" `
	IsProd          bool                  `json:"isProd" description:"Check System Env Is Prod, true|false" `
	SupportTimeZone []string              `json:"supportTimeZone" description:"Support TimeZone List" `
	SupportCurrency []*SupportCurrency    `json:"supportCurrency" description:"Support Currency List" `
	Gateway         []*ro.GatewaySimplify `json:"gateway" description:"Support Currency List" `
}

type SupportCurrency struct {
	Currency string
	Symbol   string
	Scale    int
}
