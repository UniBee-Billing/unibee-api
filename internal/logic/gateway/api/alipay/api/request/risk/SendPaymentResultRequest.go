package risk

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseRisk "unibee/internal/logic/gateway/api/alipay/api/response/risk"
)

type SendPaymentResultRequest struct {
	ReferenceTransactionId string                        `json:"referenceTransactionId,omitempty"`
	PaymentStatus          string                        `json:"paymentStatus,omitempty"`
	AuthorizationError     *model.AuthorizationError     `json:"authorizationError,omitempty"`
	CardVerificationResult *model.CardVerificationResult `json:"cardVerificationResult,omitempty"`
	PaymentMethodProvider  string                        `json:"paymentMethodProvider,omitempty"`
}

func NewSendPaymentResultRequest() (*request.AlipayRequest, *SendPaymentResultRequest) {
	sendPaymentResultRequest := &SendPaymentResultRequest{}
	alipayRequest := request.NewAlipayRequest(sendPaymentResultRequest, model.RISK_SEND_PAYMENT_RESULT_PATH, &responseRisk.SendPaymentResultResponse{})
	return alipayRequest, sendPaymentResultRequest
}
