package dispute

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseDispute "unibee/internal/logic/gateway/api/alipay/api/response/dispute"
)

type AlipaySupplyDefenseDocumentRequest struct {
	DisputeId       string `json:"disputeId,omitempty"`
	DisputeEvidence string `json:"disputeEvidence,omitempty"`
}

func NewAlipaySupplyDefenseDocumentRequest() (*request.AlipayRequest, *AlipaySupplyDefenseDocumentRequest) {
	alipaySupplyDefenseDocumentRequest := &AlipaySupplyDefenseDocumentRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySupplyDefenseDocumentRequest, model.SUPPLY_DEFENCE_DOC_PATH, &responseDispute.AlipaySupplyDefenseDocumentResponse{})
	return alipayRequest, alipaySupplyDefenseDocumentRequest
}
