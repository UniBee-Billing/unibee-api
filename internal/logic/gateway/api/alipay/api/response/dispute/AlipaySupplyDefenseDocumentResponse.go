package responseDispute

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipaySupplyDefenseDocumentResponse struct {
	response.AlipayResponse
	DisputeId             string `json:"disputeId,omitempty"`
	DisputeResolutionTime string `json:"disputeResolutionTime,omitempty"`
}
