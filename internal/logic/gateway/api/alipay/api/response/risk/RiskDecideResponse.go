package responseRisk

import "unibee/internal/logic/gateway/api/alipay/api/response"

type RiskDecideResponse struct {
	response.AlipayResponse
	Decision               string `json:"decision,omitempty"`
	AuthenticationDecision string `json:"authenticationDecision,omitempty"`
}
