package responseDispute

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayAcceptDisputeResponse struct {
	response.AlipayResponse
	DisputeId             string `json:"disputeId,omitempty"`
	DisputeResolutionTime string `json:"disputeResolutionTime,omitempty"`
}
