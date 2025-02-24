package marketplace

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseMarketplace "unibee/internal/logic/gateway/api/alipay/api/response/marketplace"
)

type AlipayRegisterRequest struct {
	RegistrationRequestId string                  `json:"registrationRequestId,omitempty"`
	SettlementInfos       []*model.SettlementInfo `json:"settlementInfos,omitempty"`
	MerchantInfo          *model.MerchantInfo     `json:"merchantInfo,omitempty"`
	PaymentMethods        []*model.PaymentMethod  `json:"paymentMethods,omitempty"`
}

func NewAlipayRegisterRequest() (*request.AlipayRequest, *AlipayRegisterRequest) {
	alipayRegisterRequest := &AlipayRegisterRequest{}
	alipayRequest := request.NewAlipayRequest(alipayRegisterRequest, model.MARKETPLACE_REGISTER_PATH, &responseMarketplace.AlipayRegisterResponse{})
	return alipayRequest, alipayRegisterRequest
}
