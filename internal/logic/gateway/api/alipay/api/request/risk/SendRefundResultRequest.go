package risk

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseRisk "unibee/internal/logic/gateway/api/alipay/api/response/risk"
)

type SendRefundResultRequest struct {
	ReferenceTransactionId string                `json:"referenceTransactionId,omitempty"`
	ReferenceRefundId      string                `json:"referenceRefundId,omitempty"`
	ActualRefundAmount     *model.Amount         `json:"actualRefundAmount,omitempty"`
	RefundRecords          []*model.RefundRecord `json:"refundRecords,omitempty"`
}

func NewSendRefundResultRequest() (*request.AlipayRequest, *SendRefundResultRequest) {
	sendRefundResultRequest := &SendRefundResultRequest{}
	alipayRequest := request.NewAlipayRequest(sendRefundResultRequest, model.RISK_SEND_REFUND_RESULT_PATH, &responseRisk.SendRefundResultResponse{})
	return alipayRequest, sendRefundResultRequest
}
