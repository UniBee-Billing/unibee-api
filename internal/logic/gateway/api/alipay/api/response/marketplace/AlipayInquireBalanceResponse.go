package responseMarketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayInquireBalanceResponse struct {
	response.AlipayResponse
	AccountBalances []*model.AccountBalance `json:"accountBalances,omitempty"`
}
