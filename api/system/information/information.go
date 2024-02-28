package information

import (
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/logic/gateway/ro"
)

type MerchantInformationReq struct {
	g.Meta `path:"/merchant_information" tags:"System-Information-Controller" method:"post" summary:"Get Merchant System Information"`
}

type MerchantInformationRes struct {
	Env             string             `description:"System Env, em: daily|stage|local|prod" `
	IsProd          bool               `description:"Check System Env Is Prod, true|false" `
	SupportTimeZone []string           `description:"Support TimeZone List" `
	SupportCurrency []*SupportCurrency `description:"Support Currency List" `
	Gateway         []*ro.GatewaySimplify
}

type SupportCurrency struct {
	Currency string
	Symbol   string
	Scale    int
}
