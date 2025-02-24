package vaulting

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseVaulting "unibee/internal/logic/gateway/api/alipay/api/response/vaulting"
)

type AlipayVaultingQueryRequest struct {
	VaultingRequestId string `json:"vaultingRequestId,omitempty"`
}

func NewAlipayVaultingQueryRequest() (*request.AlipayRequest, *AlipayVaultingQueryRequest) {
	alipayVaultingPaymentMethodRequest := &AlipayVaultingQueryRequest{}
	alipayRequest := request.NewAlipayRequest(alipayVaultingPaymentMethodRequest, model.INQUIRE_VAULTING_PATH, &responseVaulting.AlipayVaultingQueryResponse{})
	return alipayRequest, alipayVaultingPaymentMethodRequest
}
