package responseDispute

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayDownloadDisputeEvidenceResponse struct {
	response.AlipayResponse
	DisputeEvidence       string                          `json:"disputeEvidence,omitempty"`
	DisputeEvidenceFormat model.DisputeEvidenceFormatType `json:"disputeEvidenceFormat,omitempty"`
}
