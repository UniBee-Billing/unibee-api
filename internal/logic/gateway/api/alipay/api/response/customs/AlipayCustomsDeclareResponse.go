package responseCustoms

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayCustomsDeclareResponse struct {
	response.AlipayResponse
	CustomsPaymentId              string                    `json:"customsPaymentId"`
	CustomsOrderId                string                    `json:"customsOrderId"`
	IdentityCheckResult           model.IdentityCheckResult `json:"identityCheckResult"`
	ClearingChannel               model.ClearingChannel     `json:"clearingChannel"`
	ClearingTransactionId         string                    `json:"clearingTransactionId"`
	CustomsProviderRegistrationId string                    `json:"customsProviderRegistrationId"`
}
