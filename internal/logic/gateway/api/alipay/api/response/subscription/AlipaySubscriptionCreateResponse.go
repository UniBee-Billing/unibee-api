package responseSubscription

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipaySubscriptionCreateResponse struct {
	response.AlipayResponse
	SchemeUrl     string `json:"schemeUrl"`
	ApplinkUrl    string `json:"applinkUrl"`
	NormalUrl     string `json:"normalUrl"`
	AppIdentifier string `json:"appIdentifier"`
}
