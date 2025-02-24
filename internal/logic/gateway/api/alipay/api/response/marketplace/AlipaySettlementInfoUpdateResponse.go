package responseMarketplace

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipaySettlementInfoUpdateResponse struct {
	response.AlipayResponse
	UpdateStatus string `json:"updateStatus,omitempty"`
}
