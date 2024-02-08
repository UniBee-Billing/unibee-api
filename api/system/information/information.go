package information

import (
	"github.com/gogf/gf/v2/frame/g"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

type MerchantInformationReq struct {
	g.Meta `path:"/merchant_information" tags:"System-Information-Controller" method:"post" summary:"Get Merchant System Information"`
}

type MerchantInformationRes struct {
	SupportTimeZone []string
	SupportCurrency []*SupportCurrency
	MerchantId      int64
	MerchantInfo    *entity.MerchantInfo
}

type SupportCurrency struct {
	Currency string
	Symbol   string
	Scale    int
}
