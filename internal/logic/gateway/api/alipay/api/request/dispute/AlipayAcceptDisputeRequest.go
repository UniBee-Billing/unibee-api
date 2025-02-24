package dispute

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseDispute "unibee/internal/logic/gateway/api/alipay/api/response/dispute"
)

type AlipayAcceptDisputeRequest struct {
	DisputeId string `json:"disputeId,omitempty"`
}

func NewAlipayAcceptDisputeRequest() (*request.AlipayRequest, *AlipayAcceptDisputeRequest) {
	alipayAcceptDisputeRequest := &AlipayAcceptDisputeRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAcceptDisputeRequest, model.ACCEPT_DISPUTE_PATH, &responseDispute.AlipayAcceptDisputeResponse{})
	return alipayRequest, alipayAcceptDisputeRequest
}
