package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipayInquireBalanceRequest struct {
	ReferenceMerchantId string `json:"referenceMerchantId,omitempty"`
}

func NewAlipayInquireBalanceRequest() (*request.AlipayRequest, *AlipayInquireBalanceRequest) {
	alipayInquireBalanceRequest := &AlipayInquireBalanceRequest{}
	alipayRequest := request.NewAlipayRequest(alipayInquireBalanceRequest, model.MARKETPLACE_INQUIREBALANCE_PATH, &responseMarketplace.AlipayInquireBalanceResponse{})
	return alipayRequest, alipayInquireBalanceRequest
}
