package responseMarketplace

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipaySettleResponse struct {
	response.AlipayResponse
	SettlementRequestId string `json:"settlementRequestId,omitempty"`
	SettlementId        string `json:"settlementId,omitempty"`
}
