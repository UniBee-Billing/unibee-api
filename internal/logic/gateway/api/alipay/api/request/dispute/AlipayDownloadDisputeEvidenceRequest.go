package dispute

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseDispute "unibee/internal/logic/gateway/api/alipay/api/response/dispute"
)

type AlipayDownloadDisputeEvidenceRequest struct {
	DisputeId           string                    `json:"disputeId,omitempty"`
	DisputeEvidenceType model.DisputeEvidenceType `json:"disputeEvidenceType,omitempty"`
}

func NewAlipayDownloadDisputeEvidenceRequest() (*request.AlipayRequest, *AlipayDownloadDisputeEvidenceRequest) {
	alipayDownloadDisputeEvidenceRequest := &AlipayDownloadDisputeEvidenceRequest{}
	alipayRequest := request.NewAlipayRequest(alipayDownloadDisputeEvidenceRequest, model.DOWNLOAD_DISPUTE_EVIDENCE_PATH, &responseDispute.AlipayDownloadDisputeEvidenceResponse{})
	return alipayRequest, alipayDownloadDisputeEvidenceRequest
}
