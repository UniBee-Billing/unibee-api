package responseVaulting

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayVaultingSessionResponse struct {
	response.AlipayResponse
	VaultingSessionData       string `json:"vaultingSessionData,omitempty"`
	VaultingSessionId         string `json:"vaultingSessionId,omitempty"`
	VaultingSessionExpiryTime string `json:"vaultingSessionExpiryTime,omitempty"`
}
