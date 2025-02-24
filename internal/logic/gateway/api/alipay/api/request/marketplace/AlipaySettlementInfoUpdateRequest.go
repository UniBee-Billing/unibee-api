package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipaySettlementInfoUpdateRequest struct {
	UpdateRequestId       string                       `json:"updateRequestId,omitempty"`
	ReferenceMerchantId   string                       `json:"referenceMerchantId,omitempty"`
	SettlementCurrency    string                       `json:"settlementCurrency,omitempty"`
	SettlementBankAccount *model.SettlementBankAccount `json:"settlementBankAccount,omitempty"`
}

func NewAlipaySettlementInfoUpdateRequest() (*request.AlipayRequest, *AlipaySettlementInfoUpdateRequest) {
	alipaySettlementInfoUpdateRequest := &AlipaySettlementInfoUpdateRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySettlementInfoUpdateRequest, model.MARKETPLACE_SETTLEMENTINFO_UPDATE_PATH, &responseMarketplace.AlipaySettlementInfoUpdateResponse{})
	return alipayRequest, alipaySettlementInfoUpdateRequest
}
