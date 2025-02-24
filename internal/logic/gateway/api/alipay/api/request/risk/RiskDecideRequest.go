package risk

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseRisk "unibee/internal/logic/gateway/api/alipay/api/response/risk"
)

type RiskDecideRequest struct {
	ReferenceTransactionId string                   `json:"referenceTransactionId,omitempty"`
	AuthorizationPhase     model.AuthorizationPhase `json:"authorizationPhase,omitempty"`
	Orders                 []*model.Order           `json:"orders,omitempty"`
	Buyer                  *model.Buyer             `json:"buyer,omitempty"`
	ActualPaymentAmount    *model.Amount            `json:"actualPaymentAmount,omitempty"`
	PaymentDetails         []*model.PaymentDetail   `json:"paymentDetails,omitempty"`
	DiscountAmount         *model.Amount            `json:"discountAmount,omitempty"`
	EnvInfo                *model.Env               `json:"envInfo,omitempty"`
}

func NewRiskDecideRequest() (*request.AlipayRequest, *RiskDecideRequest) {
	riskDecideRequest := &RiskDecideRequest{}
	alipayRequest := request.NewAlipayRequest(riskDecideRequest, model.RISK_DECIDE_PATH, &responseRisk.RiskDecideResponse{})
	return alipayRequest, riskDecideRequest
}
