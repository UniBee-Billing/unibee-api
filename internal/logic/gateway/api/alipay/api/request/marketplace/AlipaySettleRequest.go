package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipaySettleRequest struct {
	SettlementRequestId string                    `json:"settlementRequestId,omitempty"`
	PaymentId           string                    `json:"paymentId,omitempty"`
	SettlementDetails   []*model.SettlementDetail `json:"settlementDetails,omitempty"`
}

func NewAlipaySettleRequest() (*request.AlipayRequest, *AlipaySettleRequest) {
	alipaySettleRequest := &AlipaySettleRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySettleRequest, model.MARKETPLACE_SETTLE_PATH, &responseMarketplace.AlipaySettleResponse{})
	return alipayRequest, alipaySettleRequest
}
