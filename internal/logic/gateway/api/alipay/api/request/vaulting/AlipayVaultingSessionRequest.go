package vaulting

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseVaulting "unibee/internal/logic/gateway/api/alipay/api/response/vaulting"
)

type AlipayVaultingSessionRequest struct {
	PaymentMethodType       string `json:"paymentMethodType,omitempty"`
	VaultingRequestId       string `json:"vaultingRequestId,omitempty"`
	VaultingNotificationUrl string `json:"vaultingNotificationUrl,omitempty"`
	RedirectUrl             string `json:"redirectUrl,omitempty"`
	MerchantRegion          string `json:"merchantRegion,omitempty"`
	Is3DSAuthentication     bool   `json:"is3DSAuthentication,omitempty"`
}

func NewAlipayVaultingSessionRequest() (*request.AlipayRequest, *AlipayVaultingSessionRequest) {
	alipayVaultingSessionRequest := &AlipayVaultingSessionRequest{}
	alipayRequest := request.NewAlipayRequest(alipayVaultingSessionRequest, model.CREATE_VAULTING_SESSION_PATH, &responseVaulting.AlipayVaultingSessionResponse{})
	return alipayRequest, alipayVaultingSessionRequest
}
