package information

import "github.com/gogf/gf/v2/frame/g"

type MerchantInformationReq struct {
	g.Meta     `path:"/merchant_information" tags:"System-Information-Controller" method:"post" summary:"Get Merchant System Information"`
	MerchantId string `p:"merchantId" dc:"MerchantId" v:"required#请输入MerchantId"`
}

type MerchantInformationRes struct {
	SupportTimeZone []string
	SupportCurrency []*SupportCurrency
}

type SupportCurrency struct {
	Currency string
	Symbol   string
	Scale    int
}
